package main

import (
	// mongo drivers
	_ "github.com/go-sql-driver/mysql"
	// web server
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware/cors"
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
	app.Post("/signup", signupRoute)
	// login route
	app.Post("/login", loginRoute)

	// create a new article / post
	app.Post("/new-post", newpostRoute)

	// get info for one post
	app.Get("/get-post/:postid", getpostRoute)

	// get articles for home page
	app.Get("/articles", getarticlesRoute)

	// get articles for a specific user
	app.Get("/user-articles/:username", userarticlesRoute)

	// remove an article
	app.Delete("/remove-article", removearticleRoute)


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