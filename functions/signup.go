package functions

import (
	"Forum/data"
	"database/sql"
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Configurer le gestionnaire pour le formulaire d'inscription
func Signup(w http.ResponseWriter, r *http.Request) {
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
		// Afficher le formulaire d'inscription
		tmpl := template.Must(template.ParseFiles("template/signup.html"))
		tmpl.Execute(w, nil)
		return
	}

	// Récupérer les données du formulaire
	form := data.SignupForm{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Vérifier si l'utilisateur existe déjà
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", form.Username).Scan(&count)
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
