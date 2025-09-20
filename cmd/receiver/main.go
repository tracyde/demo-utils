package main

// @title Translator Receiver API
// @version 1.0.0
// @description This is a simple API, that receives a JSON object from another application

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @servers.url http://localhost:8081

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/tracyde/demo-utils/api/resource/health"
	"github.com/tracyde/demo-utils/api/resource/object"
	"github.com/tracyde/demo-utils/config"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func handlerFuncObject(w http.ResponseWriter, r *http.Request) {
	var o object.Object

	// Decode JSON from request body
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to decode JSON")
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"object": o,
	}).Info("Received object")
}

func main() {
	config := config.New()

	log.Info("Starting receiver")
	log.WithFields(log.Fields{
		"port":         config.Server.Port,
		"inputLogPath": config.Server.InputLogPath,
	}).Info("Configuration complete")

	http.HandleFunc("/livez", health.Read)
	http.HandleFunc("/ingest", handlerFuncObject)
	log.WithFields(log.Fields{"port": config.Server.Port}).Info("Listening")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Server.Port), nil))
}
