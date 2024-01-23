package functions

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func UbanUser(w http.ResponseWriter, r *http.Request) {
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
	type BanRequest struct {
		UserID string `json:"userID"`
	}

	var banRequest BanRequest
	err = json.NewDecoder(r.Body).Decode(&banRequest)
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

	// Bannir l'utilisateur en mettant à jour la colonne "isbanned" dans la table "users"
	_, err = db.Exec("UPDATE users SET is_banned = FALSE WHERE id = ?", banRequest.UserID)
	if err != nil {
		http.Error(w, "Une erreur est survenue lors du bannissement de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page d'administration après le bannissement
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
