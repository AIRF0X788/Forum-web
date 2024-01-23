package functions

import (
	"html/template"
	"log"
	"net/http"
)

func Plus(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/plus.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.ExecuteTemplate(w, "plus.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
