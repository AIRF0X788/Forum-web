package functions

import (
	"Forum/data"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func Admin(w http.ResponseWriter, r *http.Request) {
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

	var users []data.User
	for rows.Next() {
		var user data.User
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
