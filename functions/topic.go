package functions

import (
	"Forum/data"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Topic(w http.ResponseWriter, r *http.Request) {
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

	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // Taile max 10 Mo
	if err != nil {
		log.Fatal(err)
	}

	topicForm := data.TopicForm{
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

	var user data.User
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
