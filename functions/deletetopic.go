package functions

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
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
