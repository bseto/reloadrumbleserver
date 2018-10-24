package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int `yaml:"port"`
}

var hub Hub

func loadConfig() (Config, error) {
	configFilename := "config.yml"
	configFile, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), configFilename))
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}
	log.Println("Loaded configuration:", config)
	return config, nil
}

func setUpExitSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,  // terminal interrupt (ctrl-c)
		syscall.SIGQUIT, // terminal quit (ctrl-\)
		syscall.SIGTERM, // termination
	)

	log.Println("Listening for exit signals, hit [CTRL+C] to quit")

	go func() {
		signal := <-c
		log.Println("Got interrupt signal:", signal)
		log.Println("Shutting down Reload Rumble server...")
		os.Exit(0)
	}()
}

func main() {
	// Setup configuration
	config, err := loadConfig()
	if err != nil {
		log.Println(err)
		log.Println("Failed to load config.yml file")
		os.Exit(1)
	}

	hub = newHub()

	go hub.run()

	// Capture (keyboard) interrupt signals for exit
	setUpExitSignals()

	// Serve
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(&hub, w, r)
	})

	log.Println("Listening on port:", config.Port)
	err = http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
