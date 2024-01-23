package functions

import (
	"Forum/data"
	"database/sql"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"time"
)

type ForumData struct {
	User           data.User
	Topics         []data.Topic
	CreatedAt      time.Time
	Username       string
	ProfilePicture string
	IsAdmin        bool
}

// Configurer le gestionnaire pour la page d'accueil
func Index(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if db != nil {
		CreateTable(db)
	} else {
		log.Fatal("Erreur de connexion à la base de données")
	}

	// Vérifier si l'utilisateur est connecté
	session, err := r.Cookie("user_id")
	if err != nil {
		// Rediriger vers la page de connexion si l'utilisateur n'est pas connecté
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var user data.User
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

	var topics []data.Topic
	for rows.Next() {
		var topic data.Topic
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

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Fatal(err)
	}
}
