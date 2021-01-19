package main

import (
	// mongo drivers
	_ "github.com/go-sql-driver/mysql"
	// web server
	"github.com/gofiber/fiber"
	// bcrypt
	"golang.org/x/crypto/bcrypt"
	// standard libraries
	"log"
	"math/rand"
	"fmt"
	"time"
	"strconv"
	"mime/multipart"
	"strings"
)


// signup handler
func signupRoute(c *fiber.Ctx) error {
	// get data
	signUpData := new(signupBody)

	if err := c.BodyParser(signUpData); err != nil {
		return c.SendString("error")
	}
	// create unqiue id
	var userId string
	userId = generateUserId()
	// create hashed password
	password := []byte(signUpData.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return c.SendString("error")
	}

	// insert into database
	insertion, err := db.Prepare("INSERT INTO users (email, username, passcode, id, dateSignedUp) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return c.SendString("error")
	}

	insertion.Exec(signUpData.Email, signUpData.Username, string(hashedPassword), userId, int32(time.Now().Unix()))

	return_data := make(map[string]string)
	return_data["email"] = signUpData.Email
	return_data["username"] = signUpData.Username
	return_data["id"] = userId

	return c.JSON(return_data)
}

// login handler
func loginRoute(c *fiber.Ctx) error {
	// get data
	loginData := new(loginBody)

	if err := c.BodyParser(loginData); err != nil {
		return c.SendString("error")
	}
	// get from databse
	rows, err := db.Query("SELECT email, username, passcode, id FROM users WHERE email=? LIMIT 1", loginData.Email)
	userData := make(map[string]string)
	for rows.Next() {
		var email, username, passcode, id string
		err = rows.Scan(&email, &username, &passcode, &id)
		if err != nil {
			return c.SendString(`error`)
		}
		userData["email"] = email
		userData["username"] = username
		userData["password"] = passcode
		userData["id"] = id
		break
	}
	if _, ok := userData["email"]; ok {
	} else {
		return c.SendString("bad email")
	}
	// compare user entered password with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(userData["password"]), []byte(loginData.Password))
	if err != nil {
		return c.SendString(`bad password`)
	}
	return_data := make(map[string]string)
	return_data = userData
	delete(return_data, "password")
	return c.JSON(return_data)
}



// new post handler
func newpostRoute(c *fiber.Ctx) error {
	userToken := c.FormValue("userToken")
	htmlContent := c.FormValue("html")
	postTitle := c.FormValue("title")
	imageFiles := []*multipart.FileHeader{}

	// retreive files
	for i:=0;i<200;i++ {
		key := "files" + strconv.Itoa(i)
		file, err := c.FormFile(key)
		if err == nil && file != nil {
			imageFiles = append(imageFiles, file)
		} else {
			break
		}
	}

	// create unqiue post id
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	bb := make([]rune, 5)
	for i := range bb {
		bb[i] = letters[rand.Intn(len(letters))]
	}
	postId := string(bb) + strconv.Itoa(int(time.Now().Unix()))

	// save image files
	for idx, x := range imageFiles {
		letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		b := make([]rune, 12)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}

		imageId := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond))) + string(b)
		imageType := strings.Split(x.Header["Content-Type"][0], "/")[1]
		// make sure good file type
		if includes(allowedImageTypes, imageType) == false {
			return c.SendString("error")
		} 

		// upload to blobs
		file, err := x.Open()
		if err != nil {
			log.Fatal(err)
		}
		//err = uploadBlob(file, azureContainerUrl, imageId+`.`+imageType)
		err = uploadBlob(file, imageId+`.`+imageType)
		if err != nil {
			return c.SendString("error")
		}
		//c.SaveFile(x, "./static/" + imageId + `.` + imageType)
		// add to database
		insertion, err := db.Prepare("INSERT INTO images (id, filetype, postid, orderNumber) VALUES (?, ?, ?, ?)")
		if err != nil {
			return c.SendString("error")
		}
		insertion.Exec(imageId, imageType, postId, idx)
	}
	// add post to database
	insert, err := db.Prepare("INSERT INTO posts (id, html, creationTime, userId, title) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return c.SendString("error")
	}
	insert.Exec(postId, htmlContent, time.Now().Unix(), userToken, postTitle)
	return_data := make(map[string]string)
	return_data["id"] = postId
	return c.JSON(return_data)
}

func getpostRoute(c *fiber.Ctx) error {
	postId := c.Params("postid")
	// get post data using a join
	rows, err := db.Query("SELECT users.username, posts.creationTime, posts.html, posts.title FROM users, posts WHERE users.id = posts.userId and posts.id = ?", postId)
	if err != nil {
		fmt.Println(err)
		return c.SendString(`error`)
	}

	var author, html, title string
	var creationTime int
	// get from post data from sql statment
	for rows.Next() {
		err = rows.Scan(&author, &creationTime, &html, &title)
		if err != nil {
			return c.SendString(`error`)
		}
		break
	}


	// get images from database
	rows, err = db.Query("SELECT id, filetype, postid, orderNumber FROM images WHERE postid = ? ORDER BY orderNumber ASC", postId)
	if err != nil {
		fmt.Println(err)
		return c.SendString("error")
	}
	var postImages []postImage
	for rows.Next() {
		imgData := new(postImage)
		err = rows.Scan(&imgData.Id, &imgData.Filetype, &imgData.Postid, &imgData.OrderNumber)
		if err != nil {
			fmt.Println(err)
			return c.SendString("error")
		}
		postImages = append(postImages, *imgData)
	}


	postData := postData{
		Author: author,
		CreationTime: creationTime,
		Html: html,
		Title: title,
		Images: postImages,
		ImageBaseUrl: imageBaseUrl,
	}

	// encode to json string
	jsonBytes, err := JSONMarshal(postData)
	if err != nil {
		fmt.Println(err)
		return c.SendString(`error`)
	}
	// return content
	c.Set("Content-Type", "application/json")
	return c.SendString(string(jsonBytes))
}

func getarticlesRoute(c *fiber.Ctx) error {
	var sqlQuery string
	// make database query for posts with images
	sqlQuery = `
		SELECT posts.title, posts.id as postId, posts.creationTime, images.id as imageId, images.filetype, users.username
		FROM posts, images, users
		WHERE users.id = posts.userId
		and images.postid = posts.id
		and images.orderNumber = 0
		ORDER BY posts.creationTime Desc
		LIMIT 30
	`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return c.SendString("error")
	}
	
	postsSlice := []postPreview{}

	for rows.Next() {
		postData := new(postPreview)
		err = rows.Scan(&postData.Title, &postData.PostId, &postData.CreationTime, &postData.ImageId, &postData.Filetype, &postData.Username)
		if err != nil {
			return c.SendString("error")
		}
		postsSlice = append(postsSlice, *postData)
	}

	// query for posts without images
	sqlQuery = `
		SELECT posts.title, posts.id AS postId, posts.creationTime, users.username
		FROM posts, users
		WHERE posts.userId = users.id
		AND posts.id NOT IN (SELECT postid FROM images)
		order by posts.creationTime DESC
		LIMIT 10;
	`
	rows, err = db.Query(sqlQuery)
	if err != nil {
		return c.SendString("error")
	}
	for rows.Next() {
		postData := new(postPreview)
		err = rows.Scan(&postData.Title, &postData.PostId, &postData.CreationTime, &postData.Username)
		if err != nil {
			return c.SendString("error")
		}
		postsSlice = append(postsSlice, *postData)
	}

	// sort posts slice by recently added
	postsSlice = sortArticlePreview(postsSlice)


	return_data := articlesReturn{
		ImageBaseUrl: imageBaseUrl,
		Articles: postsSlice,
	}

	return c.JSON(return_data)
}

func userarticlesRoute(c *fiber.Ctx) error {
	username := c.Params("username")

	postsSlice := []postPreview{}

	// query database for all posts with images
	sqlStatement := `
		SELECT posts.title, posts.id as postId, posts.creationTime, images.id as imageId, images.filetype
		FROM posts, images, users
		WHERE users.username = ?
		and posts.userId = users.id
		and images.postid = posts.id
		and images.orderNumber = 0
		ORDER BY posts.creationTime Desc;
	`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		return c.SendString("error")
	}

	for rows.Next() {
		postData := new(postPreview)
		err = rows.Scan(&postData.Title, &postData.PostId, &postData.CreationTime, &postData.ImageId, &postData.Filetype)
		if err != nil {
			return c.SendString("error")
		}
		postsSlice = append(postsSlice, *postData)
	}

	// query database for all posts without images
	sqlStatement = `
		SELECT posts.title, posts.id AS postId, posts.creationTime
		FROM posts, users
		WHERE users.username = ?
		AND posts.userId = users.id
		AND posts.id NOT IN (SELECT postid FROM images)
		order by posts.creationTime DESC
	`

	rows, err = db.Query(sqlStatement, username)
	if err != nil {
		return c.SendString("error")
	}
	for rows.Next() {
		postData := new(postPreview)
		err = rows.Scan(&postData.Title, &postData.PostId, &postData.CreationTime)
		if err != nil {
			return c.SendString("error")
		}
		postsSlice = append(postsSlice, *postData)
	}

	// ensure data is sorted
	postsSlice = sortArticlePreview(postsSlice)
	
	return_data := articlesReturn{
		ImageBaseUrl: imageBaseUrl,
		Articles: postsSlice,
	}
	return c.JSON(return_data)
}

func removearticleRoute(c *fiber.Ctx) error {
	body := new(deleteBody)

	if err := c.BodyParser(body); err != nil {
		c.SendString("error")
	}

	// select all images for the post
	// this join will ensure that if the user supplies the wrong userid, we will delete no images
	sqlJoin := `
		SELECT images.id as imageId, images.filetype as imageFiletype
		FROM images, posts
		WHERE posts.id = images.postId
		AND posts.id = ?
		AND posts.userId = ?
	`
	rows, err := db.Query(sqlJoin, body.PostId, body.UserId)
	if err != nil{
		c.SendString("error")
	}
	// store all images in the post Images slice
	postImages := []imageData{}
	for rows.Next() {
		image := new(imageData)
		err = rows.Scan(&image.Id, &image.Filetype)
		if err != nil {
			c.SendString("error")
		}
		postImages = append(postImages, *image)
	}


	// remove from posts table and images table
	deletion, err := db.Prepare("DELETE posts, images FROM posts INNER JOIN images ON posts.id = images.postid WHERE posts.id=? AND posts.userId=?")
	if err != nil {
		return c.SendString("error")
	}
	deletion.Exec(body.PostId, body.UserId)


	// delete images from azure blobs
	for _, img := range postImages {
		filename := img.Id + `.` + img.Filetype
		deleteBlob(filename)
	}
	return c.SendString("done")

}