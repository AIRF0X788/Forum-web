package functions

import (
	"net/http"
	"time"
)

// Configurer le gestionnaire pour la route de déconnexion
func Logout(w http.ResponseWriter, r *http.Request) {
	// Effacer le cookie contenant l'ID de l'utilisateur
	cookie := http.Cookie{Name: "user_id", Value: "", Expires: time.Unix(0, 0)}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
