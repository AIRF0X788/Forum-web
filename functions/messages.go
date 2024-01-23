package functions

import (
	"Forum/data"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

func Messages(w http.ResponseWriter, r *http.Request) {
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
	var message data.Message
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

	var user data.User
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
