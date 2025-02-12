package controllers

import (
	"fmt"
	"net/http"
	"taskManager/templates"
)

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.Templates.ExecuteTemplate(w, "base.html", map[string]string{
		"Title": "Home Page",
	})
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		fmt.Println(err)
	}
}
