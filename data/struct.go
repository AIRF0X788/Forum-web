package data

import (
	"html/template"
	"time"
)

type Message struct {
	ID             string
	User           User
	Username       string
	Contenu        string
	Valeur         string
	ProfilePicture []byte
	CreatedAt      time.Time
}

type User struct {
	ID             int
	Username       string
	PasswordHash   string
	Email          string
	CreatedAt      time.Time
	ProfilePicture []byte
	IsAdmin        bool
	IsBanned       bool
}

type SignupForm struct {
	Username string
	Email    string
	Password string
}

type LoginForm struct {
	Login    string
	Password string
}

type ForumData struct {
	User           User
	Topics         []Topic
	CreatedAt      time.Time
	Username       string
	ProfilePicture string
	IsAdmin        bool
}

type Topic struct {
	ID             int
	Title          string
	Content        string
	Category       string
	CreatedAt      time.Time
	User           User
	Username       string
	Upvotes        int
	Messages       []Message
	Image          []byte
	ImageData      template.HTML
	ProfilePicture string
}

// TopicForm est la structure pour les données envoyées par le formulaire de création de topic
type TopicForm struct {
	Title    string
	Content  string
	Category string
}
