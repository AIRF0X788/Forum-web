package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	ID        string
	User      User
	Username  string
	Contenu   string
	Valeur    string
	CreatedAt time.Time
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

func main() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		profile_picture BLOB,
		is_admin INTEGER DEFAULT 0 CHECK (is_admin IN (0, 1)),
		is_banned BOOLEAN NOT NULL DEFAULT FALSE
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS topics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        category INTEGER NOT NULL,
        upvotes TEXT NOT NULL DEFAULT '0',
        image BLOB
    );
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		contenu TEXT,
		created_at TIMESTAMP NOT NULL,
		valeur TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}

	// Créer un routeur pour l'application
	router := http.NewServeMux()

	// Configurer le gestionnaire pour le formulaire d'inscription
	signupHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// Afficher le formulaire d'inscription
			tmpl := template.Must(template.ParseFiles("template/signup.html"))
			tmpl.Execute(w, nil)
			return
		}

		// Récupérer les données du formulaire
		form := SignupForm{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// Vérifier si l'utilisateur existe déjà
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", form.Username).Scan(&count)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "Nom d'utilisateur déjà pris", http.StatusBadRequest)
			return
		}

		// Vérifier si l'email existe déjà
		var emailCount int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", form.Email).Scan(&emailCount)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}
		if emailCount > 0 {
			http.Error(w, "Email déjà existant", http.StatusBadRequest)
			return
		}

		// Vérifier si l'adresse e-mail contient un "@"
		match, _ := regexp.MatchString("[^@]+@[^@]+\\.[^@]+", form.Email)
		if !match {
			http.Error(w, "Adresse e-mail invalide", http.StatusBadRequest)
			return
		}

		// Ajouter l'utilisateur à la base de données
		hashedPassword, err := hashPassword(form.Password)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		// Déterminer si l'utilisateur est le premier à s'inscrire
		isFirstUser := false
		err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			isFirstUser = true
		}

		// Insérer l'utilisateur dans la base de données
		createdAt := time.Now()
		_, err = db.Exec("INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
			form.Username, form.Email, hashedPassword, createdAt)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		// Gérer le téléchargement de la photo de profil
		profilePicture, _, err := r.FormFile("profile_picture")
		if err == nil {
			defer profilePicture.Close()

			profilePictureBytes, err := ioutil.ReadAll(profilePicture)
			if err != nil {
				log.Println(err)
			} else {
				_, err = db.Exec("UPDATE users SET profile_picture = ? WHERE username = ?", profilePictureBytes, form.Username)
				if err != nil {
					log.Println(err)
				}
			}
		}

		// Définissez le chemin d'accès à l'image de profil
		profilePicturePath := "image/admin.png"

		// Lire le contenu du fichier de l'image
		imageData, err := ioutil.ReadFile(profilePicturePath)
		if err != nil {
			log.Println(err)
			return
		}

		// Convertir les données de l'image en chaîne Base64
		imageBase64 := base64.StdEncoding.EncodeToString(imageData)

		// Si c'est le premier utilisateur, créer un compte administrateur
		if isFirstUser {
			adminPassword := "admin123"
			adminPasswordHash, err := hashPassword(adminPassword)
			if err != nil {
				log.Println(err)
			} else {
				// Insérer le nouvel utilisateur dans la base de données avec la photo de profil en tant que chaîne Base64
				_, err = db.Exec("INSERT INTO users (username, email, password_hash, created_at, is_admin, is_banned, profile_picture) VALUES (?, ?, ?, ?, ?, ?, ?)",
					"admin", "admin@admin.com", adminPasswordHash, createdAt, 1, false, imageBase64)
				if err != nil {
					log.Println(err)
				}
			}

		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	// Configurer le gestionnaire pour le formulaire de connexion
	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// Afficher le formulaire de connexion
			tmpl := template.Must(template.ParseFiles("template/login.html"))
			tmpl.Execute(w, nil)
			return
		}
		// Récupérer les données du formulaire
		form := LoginForm{
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}
		// Récupérer l'utilisateur de la base de données
		var user User
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE username = ? OR email = ?", form.Login, form.Login).Scan(&user.ID, &user.PasswordHash)
		if err != nil {
			log.Println(err)
			http.ServeFile(w, r, "template/error.html")
			return
		}

		// Vérifier le mot de passe
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(form.Password)); err != nil {
			http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusBadRequest)
			return
		}
		// Stocker l'ID de l'utilisateur dans un cookie
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{Name: "user_id", Value: fmt.Sprint(user.ID), Expires: expiration}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Configurer le gestionnaire pour la route de déconnexion
	logoutHandler := func(w http.ResponseWriter, r *http.Request) {
		// Effacer le cookie contenant l'ID de l'utilisateur
		cookie := http.Cookie{Name: "user_id", Value: "", Expires: time.Unix(0, 0)}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	// Configurer le gestionnaire pour la page d'accueil
	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		// Vérifier si l'utilisateur est connecté
		session, err := r.Cookie("user_id")
		if err != nil {
			// Rediriger vers la page de connexion si l'utilisateur n'est pas connecté
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var user User
		err = db.QueryRow("SELECT id, username, email, created_at, profile_picture, is_admin FROM users WHERE id = ?", session.Value).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.ProfilePicture, &user.IsAdmin)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		// Récupérer les topics depuis la base de données
		rows, err := db.Query("SELECT id, title, content, created_at, category FROM topics ORDER BY created_at DESC LIMIT 10")
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		userID := session.Value

		var isBanned bool
		err = db.QueryRow("SELECT is_banned FROM users WHERE id = ?", userID).Scan(&isBanned)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		if isBanned {
			http.ServeFile(w, r, "template/ban.html")
			return
		}

		var topics []Topic
		for rows.Next() {
			var topic Topic
			err := rows.Scan(&topic.ID, &topic.Title, &topic.Content, &topic.CreatedAt, &topic.Category)
			if err != nil {
				log.Println(err)
				http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
				return
			}
			topics = append(topics, topic)
		}
		if err = rows.Err(); err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		// Convertir l'image en base64
		profilePictureBase64 := base64.StdEncoding.EncodeToString(user.ProfilePicture)

		data := ForumData{
			User:           user,
			Topics:         topics,
			CreatedAt:      time.Now(),
			Username:       user.Username,
			ProfilePicture: profilePictureBase64,
			IsAdmin:        user.IsAdmin,
		}

		data.GetData(user, topics)
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			log.Fatal(err)
		}
	}

	DeleteTopicHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		userID := cookie.Value

		topicID := r.FormValue("topic_id")

		// Convertir l'ID du topic en entier
		topicIDInt, err := strconv.Atoi(topicID)
		if err != nil {
			http.Error(w, "ID de topic invalide", http.StatusBadRequest)
			return
		}

		// Vérifier si l'utilisateur est administrateur
		isAdmin := 0
		err = db.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&isAdmin)
		if err != nil {
			// Gérer l'erreur
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		if isAdmin != 1 {
			http.ServeFile(w, r, "template/404.html")
			return
		}

		// Supprimer le topic de la base de données en utilisant l'ID
		_, err = db.Exec("DELETE FROM topics WHERE id = ?", topicIDInt)
		if err != nil {
			// Gérer l'erreur
			http.Error(w, "Une erreur est survenue lors de la suppression du topic", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/topics", http.StatusSeeOther)
	}

	DeleteMessageHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		userID := cookie.Value

		messageID := r.FormValue("message_id")

		var messageIDInt int
		var err error

		messageIDInt, err = strconv.Atoi(messageID)
		if err != nil {
			http.Error(w, "ID de message invalide", http.StatusBadRequest)
			return
		}

		// Vérifier si l'utilisateur est administrateur
		isAdmin := 0
		err = db.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&isAdmin)
		if err != nil {
			// Gérer l'erreur
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		if isAdmin != 1 {
			http.ServeFile(w, r, "template/404.html")
			return
		}

		// Supprimer le topic de la base de données en utilisant l'ID
		_, err = db.Exec("DELETE FROM messages WHERE id = ?", messageIDString)
		if err != nil {
			// Gérer l'erreur
			http.Error(w, "Une erreur est survenue lors de la suppression du topic", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/topics", http.StatusSeeOther)
	}

	TopicHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // Taile max 10 Mo
		if err != nil {
			log.Fatal(err)
		}

		topicForm := TopicForm{
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
			Category: r.FormValue("category"),
		}

		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		userID := cookie.Value

		var user User
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&user.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO topics (username, title, content, created_at, category) VALUES (?, ?, ?, ?, ?)",
			user.Username, topicForm.Title, topicForm.Content, time.Now(), topicForm.Category)
		if err != nil {
			log.Fatal(err)
		}

		// dl de l'image
		file, _, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			imageBytes, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println(err)
			} else {
				var topicID int64
				err = db.QueryRow("SELECT last_insert_rowid()").Scan(&topicID)
				if err != nil {
					log.Println(err)
				} else {
					_, err = db.Exec("UPDATE topics SET image = ? WHERE id = ?", imageBytes, topicID)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

		http.Redirect(w, r, "/topics", http.StatusSeeOther)
	}

	TopicViewHandler := func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Erreur de chargement de la session", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		idStr := r.URL.Path[len("/topic/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalide topic ID", http.StatusBadRequest)
			return
		}
		if session == nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}

		// Récupérer l'utilisateur de la base de données
		var user User
		err = db.QueryRow("SELECT is_admin FROM users WHERE id = ?", session.Value).Scan(&user.IsAdmin)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		// Récupérer les données du topic à partir de la base de données
		var topic Topic
		err = db.QueryRow("SELECT id, title, content, created_at, category, username, upvotes, image FROM topics WHERE id = ?", id).Scan(&topic.ID, &topic.Title, &topic.Content, &topic.CreatedAt, &topic.Category, &topic.Username, &topic.Upvotes, &topic.Image)
		if err != nil {
			log.Println(err)
			http.Error(w, "Topic introuvable", http.StatusNotFound)
			return
		}

		// PARITE MESSAGE ---------------------------------------------------------------------------------------------------

		// Récupérer les messages associés au topic à partir de la base de données
		rows, err := db.Query("SELECT id, username,contenu, created_at,valeur FROM messages WHERE valeur = ?", "/topic/"+idStr)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var messages []Message
		for rows.Next() {
			var message Message
			err := rows.Scan(&message.ID, &message.Username, &message.Contenu, &message.CreatedAt, &message.Valeur)
			if err != nil {
				log.Println(err)
				http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
				return
			}
			messages = append(messages, message)
		}
		if err = rows.Err(); err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		// Ajouter les messages au topic
		topic.Messages = messages

		// --------------------------------------------------------------------------------------------------------------------

		if topic.Username != "" {
			user.Username = topic.Username
		} else {
			_, err = db.Exec("UPDATE topics SET username = ? WHERE id = ?", user.Username, topic.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to update topic with username", http.StatusInternalServerError)
				return
			}
		}

		action := r.FormValue("action")
		switch action {
		case "upvote":
			_, err = db.Exec("UPDATE topics SET upvotes = upvotes + 1 WHERE id = ?", id)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to upvote topic", http.StatusInternalServerError)
				return
			}
			topic.Upvotes++
		case "downvote":
			_, err = db.Exec("UPDATE topics SET upvotes = upvotes - 1 WHERE id = ?", id)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to downvote topic", http.StatusInternalServerError)
				return
			}
			topic.Upvotes--
		}

		// Prépare les données de l'image
		var imageData string
		if topic.Image != nil {
			imageData = base64.StdEncoding.EncodeToString(topic.Image)
		}

		profilePictureBase64 := base64.StdEncoding.EncodeToString(user.ProfilePicture)

		data := Topic{
			Title:          topic.Title,
			Content:        topic.Content,
			Category:       topic.Category,
			CreatedAt:      time.Now(),
			User:           user,
			Username:       user.Username,
			Upvotes:        topic.Upvotes,
			Messages:       messages,
			Image:          topic.Image,
			ImageData:      template.HTML(imageData),
			ProfilePicture: profilePictureBase64,
		}
		data.GetUser(user)
		tmpl, err := template.ParseFiles("template/topic.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.ExecuteTemplate(w, "topic.html", data)
		if err != nil {
			log.Fatal(err)
		}
	}

	plusHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/plus.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.ExecuteTemplate(w, "plus.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	updateUsernameHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// Afficher le formulaire d'inscription
			tmpl := template.Must(template.ParseFiles("template/usersrename.html"))
			tmpl.Execute(w, nil)
			return
		}
		// Vérifier si l'utilisateur est connecté
		session, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Récupérer le nouvel username depuis le formulaire
		newUsername := r.FormValue("username")

		// Mettre à jour le nom d'utilisateur dans la base de données
		_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, session.Value)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	changeProfilePictureHandler := func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID := session.Value

		// Vérifier si un fichier a été envoyé dans la demande
		file, _, err := r.FormFile("profilePicture")
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue lors de l'envoi du fichier", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		imageData, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue lors de la lecture du fichier", http.StatusInternalServerError)
			return
		}

		// Mettre à jour la photo de profil dans la bdd
		_, err = db.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", imageData, userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue lors de la mise à jour de la photo de profil", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	funnyUsernames := []string{
		"CaptainPoop",
		"MrFartyPants",
		"SirStinky",
		"PrincessPoo",
		"BaronVonStench",
		"DukeDoodoo",
		"LordOfPoop",
		"FartMaster",
		"PoopMan",
	}

	generateUsernameHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		session, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userID := session.Value

		// génere pseudo aléatoire depuis la liste
		rand.Seed(time.Now().Unix())
		randomIndex := rand.Intn(len(funnyUsernames))

		randomUsername := funnyUsernames[randomIndex]

		_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", randomUsername, userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	messageHandler := func(w http.ResponseWriter, r *http.Request) {
		// Récupérer l'URL complète de la page qui a effectué la requête
		referer := r.Referer()

		// Extraire l'URL visible à partir de l'URL complète
		url, err := url.Parse(referer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		visibleURL := url.Path

		// Récupérer les données JSON envoyées dans la requête
		var message Message
		err = json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		userID := cookie.Value

		var user User
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&user.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		createdAt := time.Now()

		// Générer un nouvel ID unique
		message.ID = uuid.New().String()

		// Insérer le message dans la base de données avec l'URL visible
		insertQuery := "INSERT INTO messages (id, username, contenu, created_at, valeur) VALUES (?, ?, ?,?,?)"
		_, err = db.Exec(insertQuery, message.ID, user.Username, message.Contenu, createdAt, visibleURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	adminHandler := func(w http.ResponseWriter, r *http.Request) {
		// Vérifier si l'utilisateur est un administrateur
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		userID := cookie.Value

		userIDCookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		connectedUserID := userIDCookie.Value

		connectedUserIDInt, err := strconv.Atoi(connectedUserID)
		if err != nil {
			http.Error(w, "Erreur de conversion de l'ID de l'utilisateur connecté", http.StatusInternalServerError)
			return
		}

		var isAdmin, isBanned bool
		err = db.QueryRow("SELECT is_admin, is_banned FROM users WHERE id = ?", userID).Scan(&isAdmin, &isBanned)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.ServeFile(w, r, "template/404.html")
			return
		}

		// Récupérer tous les utilisateurs de la base de données
		rows, err := db.Query("SELECT id, username, email, is_admin, is_banned FROM users")
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de la récupération des utilisateurs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.IsBanned)
			if err != nil {
				http.Error(w, "Une erreur est survenue lors de la récupération des utilisateurs", http.StatusInternalServerError)
				return
			}
			// Ajoutez cette vérification pour marquer l'utilisateur connecté comme banni
			if user.ID == connectedUserIDInt {
				user.IsBanned = true
			}
			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Une erreur est survenue lors de la récupération des utilisateurs", http.StatusInternalServerError)
			return
		}

		// Charger le modèle HTML de la page d'administration
		tmpl := template.Must(template.ParseFiles("template/admin.html"))

		// Passer les données des utilisateurs au modèle HTML
		err = tmpl.Execute(w, users)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de l'affichage de la page d'administration", http.StatusInternalServerError)
			return
		}
	}

	banUserHandler := func(w http.ResponseWriter, r *http.Request) {
		type BanRequest struct {
			UserID string `json:"userID"`
		}

		var banRequest BanRequest
		err := json.NewDecoder(r.Body).Decode(&banRequest)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de la lecture de la requête", http.StatusBadRequest)
			return
		}

		userIDCookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		connectedUserID := userIDCookie.Value

		var isAdmin bool
		err = db.QueryRow("SELECT is_admin FROM users WHERE id = ?", connectedUserID).Scan(&isAdmin)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.ServeFile(w, r, "template/404.html")
			return
		}

		if banRequest.UserID == "" {
			http.Error(w, "Identifiant d'utilisateur manquant", http.StatusBadRequest)
			return
		}

		// Bannir l'utilisateur en mettant à jour la colonne "is_banned" dans la table "users"
		_, err = db.Exec("UPDATE users SET is_banned = TRUE WHERE id = ?", banRequest.UserID)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors du bannissement de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Rediriger vers la page d'administration après le bannissement
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	unbanUserHandler := func(w http.ResponseWriter, r *http.Request) {
		type BanRequest struct {
			UserID string `json:"userID"`
		}

		var banRequest BanRequest
		err := json.NewDecoder(r.Body).Decode(&banRequest)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de la lecture de la requête", http.StatusBadRequest)
			return
		}

		userIDCookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}
		connectedUserID := userIDCookie.Value

		var isAdmin bool
		err = db.QueryRow("SELECT is_admin FROM users WHERE id = ?", connectedUserID).Scan(&isAdmin)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.ServeFile(w, r, "template/404.html")
			return
		}

		if banRequest.UserID == "" {
			http.Error(w, "Identifiant d'utilisateur manquant", http.StatusBadRequest)
			return
		}

		// Bannir l'utilisateur en mettant à jour la colonne "is_banned" dans la table "users"
		_, err = db.Exec("UPDATE users SET is_banned = FALSE WHERE id = ?", banRequest.UserID)
		if err != nil {
			http.Error(w, "Une erreur est survenue lors du bannissement de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Rediriger vers la page d'administration après le bannissement
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	router.HandleFunc("/image/logoutico.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/logoutico.png")
	})

	router.HandleFunc("/image/mode-sombre.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/mode-sombre.png")
	})

	router.HandleFunc("/image/wallpaper.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/wallpaper.jpg")
	})

	router.HandleFunc("/image/wallpaper1.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/wallpaper1.jpg")
	})

	router.HandleFunc("/image/engrenage.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/engrenage.png")
	})

	router.HandleFunc("/image/back.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/back.png")
	})

	router.HandleFunc("/image/admin.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/admin.png")
	})

	// Configurer le gestionnaire pour les fichiers statiques
	router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	router.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	router.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir("audio"))))

	// Configurer les routes
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/create_topic", TopicHandler)
	router.HandleFunc("/topic/", TopicViewHandler)
	router.HandleFunc("/signup", signupHandler)
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/plus", plusHandler)
	router.HandleFunc("/enregistrer-message", messageHandler)
	router.HandleFunc("/update", updateUsernameHandler)
	router.HandleFunc("/generate-username", generateUsernameHandler)
	router.HandleFunc("/change-profile-picture", changeProfilePictureHandler)
	router.HandleFunc("/delete_topic", DeleteTopicHandler)
	router.HandleFunc("/admin", adminHandler)
	router.HandleFunc("/ban", banUserHandler)
	router.HandleFunc("/unban", unbanUserHandler)
	router.HandleFunc("/delete_message", DeleteMessageHandler)

	// Démarrer le serveur HTTP
	log.Println("Serveur démarré sur http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func (fd *ForumData) GetData(user User, topics []Topic) {
	fd.User = user
	fd.Topics = topics
}

func (fd *Topic) GetUser(user User) {
	fd.User = user
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
