package functions

import (
	"Forum/data"
	"database/sql"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func TopicView(w http.ResponseWriter, r *http.Request) {
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
	var user data.User
	err = db.QueryRow("SELECT is_admin FROM users WHERE id = ?", session.Value).Scan(&user.IsAdmin)
	if err != nil {
		log.Println(err)
		http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
		return
	}

	// Récupérer les données du topic à partir de la base de données
	var topic data.Topic
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

	var messages []data.Message
	for rows.Next() {
		var message data.Message
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

	data := data.Topic{
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

	tmpl, err := template.ParseFiles("template/topic.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.ExecuteTemplate(w, "topic.html", data)
	if err != nil {
		log.Fatal(err)
	}
}
