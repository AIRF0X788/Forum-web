package functions

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func GenerateUsername(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	funnyUsernames := []string{
		"CaptainPoop",
		"MrFartyPants",
		"SirStinky",
		"PrincessPoo",
		"BaronVonStench",
		"DukeDoodoo",
		"LordOfPoop",
		"FartMaster",
		"PoopMan",
	}

	session, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID := session.Value

	// génere pseudo aléatoire depuis la liste
	rand.Seed(time.Now().Unix())
	randomIndex := rand.Intn(len(funnyUsernames))

	randomUsername := funnyUsernames[randomIndex]

	_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", randomUsername, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Une erreur est survenue", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
