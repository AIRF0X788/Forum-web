package functions

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func UpdateUsername(w http.ResponseWriter, r *http.Request) {
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
		tmpl := template.Must(template.ParseFiles("template/usersrename.html"))
		tmpl.Execute(w, nil)
		return
	}
	// Vérifier si l'utilisateur est connecté
	session, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer le nouvel username depuis le formulaire
	newUsername := r.FormValue("username")

	// Mettre à jour le nom d'utilisateur dans la base de données
	_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, session.Value)
	if err != nil {
		log.Println(err)
		http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
