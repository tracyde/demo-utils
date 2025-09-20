package main

// @title Translator API
// @version 1.0.0
// @description This is a simple translator API, that pulls translations from a log file and sends them to a receiver after converting to JSON format.

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @servers.url http://localhost:8081

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jdrews/go-tailer/fswatcher"
	"github.com/jdrews/go-tailer/glob"
	"github.com/tracyde/demo-utils/api/resource/object"
	"github.com/tracyde/demo-utils/config"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	config := config.New()

	log.Info("Starting translator API")
	log.WithFields(log.Fields{
		"port":         config.Server.Port,
		"inputLogPath": config.Server.InputLogPath,
	}).Info("Configuration complete")

	// Read from log file looking for translations
	//
	path := config.Server.InputLogPath
	parsedGlob, err := glob.Parse(path)
	if err != nil {
		log.Fatalf("Failed to parse InputLogPath %s: %v", path, err)
	}

	log.Info("Initiating file tailer: ", path)
	// specify how often you want the tailer to check for updates
	pollInterval := time.Duration(500 * time.Millisecond)
	// TODO: handle invalid InputLogPath, below is supposed to catch it but does not seem to work properly
	// tailer, err := fswatcher.RunFileTailer([]glob.Glob{parsedGlob}, false, true, log.New())
	tailer, err := fswatcher.RunPollingFileTailer([]glob.Glob{parsedGlob}, false, true, pollInterval, log.New())
	if err != nil {
		log.Fatalf("Failed to start file tailer: %v", err)
	}
	defer tailer.Close()

	// listen to the go channel for captured lines
	url := "http://" + config.Server.Host + ":" + strconv.Itoa(config.Server.Port) + config.Server.Endpoint
	token := "ja -> en"
	for line := range tailer.Lines() {
		log.Debug("Captured line:", line)
		if strings.Contains(line.Line, token) {
			err := processLine(url, line)
			if err != nil {
				log.Errorf("Failed to process line: %v", err)
			}
		}
	}

}

// TODO: write logic to process translations (for now just generating fixed translation objects)
// Process translations (for now just generating fixed translation objects)
func processLine(url string, line *fswatcher.Line) error {
	log.Info("Found translation:", line.Line)

	p, err := chooseObject(line)
	if err != nil {
		return err
	}

	// TODO: add configuration to set LOG_LEVEL from environment variable
	log.Debug("Chose object: ", p)

	// TODO: write logic to send translations to receiver
	// Send translations to receiver
	err = sendObject(url, p)

	return err
}

func sendObject(url string, o object.Object) error {
	j, err := json.Marshal(o)
	if err != nil {
		return err
	}

	jsonStr := string(j)
	log.Debug("JSON object: ", jsonStr)

	log.WithFields(log.Fields{
		"url":  url,
		"json": jsonStr,
	}).Info("Sending JSON object")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{
		"status":  resp.Status,
		"headers": resp.Header,
	}).Info("Sent JSON Object")

	return nil
}

func chooseObject(line *fswatcher.Line) (object.Object, error) {
	l := strings.ToLower(line.Line)
	var o object.Object
	var err error

	switch {
	case strings.Contains(l, "viper"):
		t, err := parseTime("2025-08-29T12:45:30Z")
		if err != nil {
			return object.Object{}, err
		}
		o = object.Object{
			Name:     "viper 1-1",
			Type:     "friendly",
			Platform: "f-16c",
			Position: object.Position{
				Latitude:  34.123,
				Longitude: 132.456,
				Altitude:  24500,
				Speed:     420,
				Heading:   275,
			},
			Status:       "on task",
			Iff:          true,
			Sensor:       "",
			TrackQuality: 0,
			Confidence:   0,
			Activity:     "",
			Time:         t,
		}
	case strings.Contains(l, "unknown"):
		t, err := parseTime("2025-08-29T12:47:10Z")
		if err != nil {
			return object.Object{}, err
		}
		o = object.Object{
			Name:     "",
			Type:     "unknown",
			Platform: "aircraft",
			Position: object.Position{
				Latitude:  35.678,
				Longitude: 129.432,
				Altitude:  15200,
				Speed:     520,
				Heading:   90,
			},
			Status:       "one unknown aircraft detected",
			Iff:          false,
			Sensor:       "aesa radar",
			TrackQuality: 7,
			Confidence:   0,
			Activity:     "",
			Time:         t,
		}
	case strings.Contains(l, "emission"):
		t, err := parseTime("2025-08-29T12:49:00Z")
		if err != nil {
			return object.Object{}, err
		}
		o = object.Object{
			Name:     "emission",
			Type:     "radar",
			Platform: "sa-21",
			Position: object.Position{
				Latitude:  36.100,
				Longitude: 129.900,
				Altitude:  0,
				Speed:     0,
				Heading:   0,
			},
			Status:       "assessed threat radius 40km",
			Iff:          false,
			Sensor:       "aesa radar",
			TrackQuality: 0,
			Confidence:   0.85,
			Activity:     "tracking",
			Time:         t,
		}
	default:
		err = fmt.Errorf("Could not choose object from given line: %s", line.Line)
	}

	return o, err
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	return t, err
}
