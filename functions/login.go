package functions

import (
	"Forum/data"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Configurer le gestionnaire pour le formulaire de connexion
func Login(w http.ResponseWriter, r *http.Request) {
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

	if r.Method != http.MethodPost {
		// Afficher le formulaire de connexion
		tmpl := template.Must(template.ParseFiles("template/login.html"))
		tmpl.Execute(w, nil)
		return
	}
	// Récupérer les données du formulaire
	form := data.LoginForm{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}
	// Récupérer l'utilisateur de la base de données
	var user data.User
	err = db.QueryRow("SELECT id, password_hash FROM users WHERE username = ? OR email = ?", form.Login, form.Login).Scan(&user.ID, &user.PasswordHash)
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
