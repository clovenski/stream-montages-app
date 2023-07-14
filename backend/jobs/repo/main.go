package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clovenski/stream-montages-app/backend/jobs/repo/config"
	"github.com/clovenski/stream-montages-app/backend/jobs/repo/controllers"
	"github.com/clovenski/stream-montages-app/backend/jobs/repo/database"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/montages/jobs", controllers.GetMontageJobs).Methods("GET")
	router.HandleFunc("/montages/jobs", controllers.CreateMontageJob).Methods("POST")
	router.HandleFunc("/montages/jobs/{id}", controllers.GetMontageJobByID).Methods("GET")
	router.HandleFunc("/montages/jobs/{id}", controllers.UpdateMontageJobByID).Methods("PUT")
	router.HandleFunc("/montages/jobs/{id}", controllers.DeleteMontageJobByID).Methods("DELETE")
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
