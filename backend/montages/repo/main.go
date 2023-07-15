package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clovenski/stream-montages-app/backend/montages/repo/config"
	"github.com/clovenski/stream-montages-app/backend/montages/repo/controllers"
	"github.com/clovenski/stream-montages-app/backend/montages/repo/database"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/montages", controllers.GetMontages).Methods("GET")
	router.HandleFunc("/montages", controllers.CreateMontage).Methods("POST")
	router.HandleFunc("/montages/{id}", controllers.GetMontageByID).Methods("GET")
	router.HandleFunc("/montages/{id}", controllers.UpdateMontageByID).Methods("PUT")
	router.HandleFunc("/montages/{id}", controllers.DeleteMontageByID).Methods("DELETE")
}

func main() {
	config, err := config.LoadDBConfig()
	if err != nil {
		log.Fatal("Failed to load app config", err)
	}

	database.Connect(&config)
	database.Migrate()

	router := mux.NewRouter().StrictSlash(true)

	controllers.InitController(database.DB)
	RegisterRoutes(router)

	log.Println(fmt.Sprintf("Starting Server on port %s", config.ServerPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.ServerPort), router))
}
