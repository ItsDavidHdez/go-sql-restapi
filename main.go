package main

import (
	"log"
	"net/http"

	"github.com/ItsDavidHdez/go-sql-restapi/db"
	"github.com/ItsDavidHdez/go-sql-restapi/models"
	"github.com/ItsDavidHdez/go-sql-restapi/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	db.DBConnect()

	db.DB.AutoMigrate(models.Music{})
	db.DB.AutoMigrate(models.User{})
	router := mux.NewRouter()

	const api = "/api/v1"
	// Music endpoints
	router.HandleFunc(api + "/songs", routes.GetAllMusicHander).Methods("GET")
	router.HandleFunc(api + "/music", routes.GetMusicHander).Methods("POST")
	router.HandleFunc(api + "/lyrics", routes.GetMusicSoapHandler).Methods("POST")

	// Users endpoints
	router.HandleFunc(api + "/", routes.HomeHandler)
	router.HandleFunc(api + "/users", routes.GetUsersHandler).Methods("GET")
	router.HandleFunc(api + "/users/{id}", routes.GetUserHandler).Methods("GET")
	router.HandleFunc(api + "/users/{id}", routes.DeleteUserHandler).Methods("DELETE")

	// Users auth endpoints
	router.HandleFunc(api + "/register", routes.CreateUserHandler).Methods("POST")
	router.HandleFunc(api + "/login", routes.LoginUserHandler).Methods("POST")
	
	http.ListenAndServe(":3000", router)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}