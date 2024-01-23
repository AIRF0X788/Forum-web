package main

import (
	"Forum/functions"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if db != nil {
		functions.CreateTable(db)
	} else {
		log.Fatal("Erreur de connexion à la base de données")
	}

	functions.CreateTable(db)

	// Créer un routeur pour l'application
	router := http.NewServeMux()

	router.HandleFunc("/image/logoutico.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/logoutico.png")
	})

	router.HandleFunc("/image/mode-sombre.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/mode-sombre.png")
	})

	router.HandleFunc("/image/wallpaper.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/wallpaper.jpg")
	})

	router.HandleFunc("/image/wallpaper1.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/wallpaper1.jpg")
	})

	router.HandleFunc("/image/engrenage.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/engrenage.png")
	})

	router.HandleFunc("/image/back.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/back.png")
	})

	router.HandleFunc("/image/admin.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./image/admin.png")
	})

	// Configurer le gestionnaire pour les fichiers statiques
	router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	router.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	router.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir("audio"))))

	// Configurer les routes
	router.HandleFunc("/", functions.Index)
	router.HandleFunc("/create_topic", functions.Topic)
	router.HandleFunc("/topic/", functions.TopicView)
	router.HandleFunc("/signup", functions.Signup)
	router.HandleFunc("/logout", functions.Logout)
	router.HandleFunc("/login", functions.Login)
	router.HandleFunc("/plus", functions.Plus)
	router.HandleFunc("/enregistrer-message", functions.Messages)
	router.HandleFunc("/update", functions.UpdateUsername)
	router.HandleFunc("/generate-username", functions.GenerateUsername)
	router.HandleFunc("/change-profile-picture", functions.ChangeProfilePicture)
	router.HandleFunc("/delete_topic", functions.DeleteTopic)
	router.HandleFunc("/admin", functions.Admin)
	router.HandleFunc("/ban", functions.BanUser)
	router.HandleFunc("/unban", functions.UbanUser)

	// Démarrer le serveur HTTP
	log.Println("Serveur démarré sur http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
