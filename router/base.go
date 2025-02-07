package router

import (
	"fmt"
	"log"
	"net/http"
	"taskManager/controllers"
)

func InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Define routes and their handlers
	mux.HandleFunc("/api/system", controllers.SystemHandler)
	mux.HandleFunc("/api/system-surface", controllers.ThirdPartySystemHandler)
	mux.HandleFunc("/api/process/start", controllers.StartProcessHandler)
	mux.HandleFunc("/api/process/stop", controllers.StopProcessHandler)

	return mux
}

func StartServer() {
	router := InitRoutes()

	// Start the server
	fmt.Println("Starting server on :7000")
	if err := http.ListenAndServe(":7000", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
