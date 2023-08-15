package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Host            = ""
	Port            = ""
	GoogleSheetId   = ""
	GoogleSheetName = ""
)

func Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Port, Host, GoogleSheetId, GoogleSheetName = os.Getenv("PORT"), os.Getenv("LISTEN_ADDR"), os.Getenv("GOOGLE_SHEET_ID"), os.Getenv("GOOGLE_SHEET_NAME")

	if Port == "" {
		Port = "4000"
	}

	if Host == "" {
		Host = "localhost"
	}

	if GoogleSheetId == "" {
		log.Fatal("unable to access google sheet id, check environment variables")
	}

	if GoogleSheetName == "" {
		log.Fatal("unable to access google sheet name, check environment variables")
	}

}
