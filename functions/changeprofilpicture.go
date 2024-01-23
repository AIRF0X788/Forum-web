package functions

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
)

func ChangeProfilePicture(w http.ResponseWriter, r *http.Request) {
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
