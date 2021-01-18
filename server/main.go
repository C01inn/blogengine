package main

import (
	// mongo drivers
	_ "github.com/go-sql-driver/mysql"
	// web server
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware/cors"
	// bcrypt
	"golang.org/x/crypto/bcrypt"
	// azure blob
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/azure-pipeline-go/pipeline"
	// standard libraries
	"log"
	"math/rand"
	"fmt"
	"time"
	"database/sql"
	"strconv"
	"mime/multipart"
	"strings"
	"encoding/json"
	"bytes"
	"net/url"
	"context"
	"io"
)





type signupBody struct {
	Email string
	Username string
	Password string
}

type loginBody struct {
	Email string
	Password string
}

type user struct {
	Email string
	Username string
	Password string
	Id string
	DateSignedUp int
}

type postData struct {
	Author string
	CreationTime int
	Html string
	Title string
	Images []postImage
	ImageBaseUrl string
}

type postImage struct {
	Id string
	Filetype string
	Postid string
	OrderNumber int
}

type postPreview struct {
	Title string
	PostId string
	CreationTime int
	ImageId string
	Filetype string
	Username string
}

type articlesReturn struct {
	ImageBaseUrl string
	Articles []postPreview
}

type deleteBody struct {
	UserId string
	PostId string
}

type imageData struct {
	Id string
	Filetype string
}

var imageBaseUrl, azureAccessKey, azureStorageAccount, containerName string
var azurePipeline pipeline.Pipeline
var azureContainerUrl azblob.ContainerURL
var db *sql.DB
var allowedImageTypes []string

func main() {
	// global variables
	imageBaseUrl = `https://blogengine.blob.core.windows.net/images/`
	allowedImageTypes = []string{"jpg", "jpeg", "png", "gif"}

	// connect to azure blobs pipeline
	azureAccessKey = "vRE++Upd7b5xmx9lrQ+rN1h3mmfiv+NhUQT7oVmfN3D54zcox7NNbB3QsgY9h3yqZF23K+qzkVF2K76dPQKR9A=="
	azureStorageAccount = "blogengine"
	containerName = "images"

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(azureStorageAccount, azureAccessKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	azurePipeline = azblob.NewPipeline(credential, azblob.PipelineOptions{})

	tempUrl, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", azureStorageAccount, containerName))

	azureContainerUrl = azblob.NewContainerURL(*tempUrl, azurePipeline)

	/******************************************************************/

	// connect to database
	db, err = sql.Open("mysql", "root:8268Wrenfield@tcp(127.0.0.1:1800)/blog")
	checkErr(err)

	/**************************************************************************************/

	app := fiber.New()
	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	// static files
	app.Static("/static", "./static")


	// index page
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	// sign up page
	app.Post("/signup", func(c *fiber.Ctx) error {
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
	})
	// login route
	app.Post("/login", func(c *fiber.Ctx) error {
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
	})

	app.Post("/new-post", func(c *fiber.Ctx) error {
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
	})




	// get info for one post
	app.Get("/get-post/:postid", func(c *fiber.Ctx) error {
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
	})

	// get articles for home page
	
	app.Get("/articles", func(c *fiber.Ctx) error {
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
	})

	app.Get("/user-articles/:username", func(c *fiber.Ctx) error {
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
	})

	app.Delete("/remove-article", func(c *fiber.Ctx) error {

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

	})


	app.Listen(":3000")
}

func includes(data []string, elem string) bool {
	elem = strings.ToLower(elem)
	for _, x := range data {
		if elem == x {
			return true
		}
	}
	return false
}


func deleteBlob(blobName string) bool {

	accessConditions := azblob.BlobAccessConditions {
	}
	fullBlobUrl, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", "blogengine", "images", blobName))

	blob := azblob.NewBlobURL(*fullBlobUrl, azurePipeline)

	blob.Delete(context.Background(), azblob.DeleteSnapshotsOptionNone, accessConditions)
	
	return true
}

func uploadBlob(file multipart.File, filename string) error {
	blobURL := azureContainerUrl.NewBlockBlobURL(filename)

	// create buffer for file
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return err
	} 

	_, err := azblob.UploadBufferToBlockBlob(context.Background(), buf.Bytes(), blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return err
	}
	return nil
}

// create a universal unqiue id for userids 
// uses random characters with datetime
func generateUserId() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 7)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
	}
	firstPart := string(b)
	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	a := make([]rune, 7)
	for i := range a {
		a[i] = letters[rand.Intn(len(letters))]
	}
	lastPart := string(a)
	return firstPart + timeStamp + lastPart
}


func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func JSONMarshal(t interface{}) ([]byte, error) {
    buffer := &bytes.Buffer{}
    encoder := json.NewEncoder(buffer)
    encoder.SetEscapeHTML(false)
    err := encoder.Encode(t)
    return buffer.Bytes(), err
}
// sorts slice of article preview type
// uses quicksort algortmn
func sortArticlePreview(arr []postPreview) []postPreview {
	if len(arr) <= 1 {
		return arr
	}
	pivot := arr[int(len(arr) / 2)]

	bigger := []postPreview{}
	smaller := []postPreview{}
	same := []postPreview{}

	for _, x := range arr {
		if x.CreationTime < pivot.CreationTime {
			bigger = append(bigger, x)
		} else if x.CreationTime > pivot.CreationTime {
			smaller = append(smaller, x)
		} else {
			same = append(same, x)
		}
	}
	return append(sortArticlePreview(smaller), append(same, sortArticlePreview(bigger)...)...)
}