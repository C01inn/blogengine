package main


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