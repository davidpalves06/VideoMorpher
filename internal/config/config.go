package config

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/davidpalves06/GodeoEffects/internal/handlers"
	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

var ApplicationConfig Configuration = Configuration{
	ServerConfig: ServerConfig{
		Host: "",
	},
	LogConfig: LogConfig{
		Level: 1,
	},
	UploadDir: "./uploads/",
}

type Configuration struct {
	ServerConfig ServerConfig `json:"server"`
	LogConfig    LogConfig    `json:"log"`
	UploadDir    string       `json:"uploadDir"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type LogConfig struct {
	Level int `json:"level"`
}

func LoadConfigurations() {
	var configFileName = *flag.String("configFile", "./config.json", "The file used to load configurations")
	flag.Parse()

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatalln("Error loading config File")
	}

	conf, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatalln("Error reading contents from config file")
	}

	err = json.Unmarshal(conf, &ApplicationConfig)

	if err != nil {
		log.Fatalln("Error parsing JSON configuration")
	}

	log.Printf("Configurations Loaded\n")

	logger.InitLogger(ApplicationConfig.LogConfig.Level)
	handlers.UPLOAD_DIRECTORY = ApplicationConfig.UploadDir
}
