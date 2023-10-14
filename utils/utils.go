package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Unable to load env. Err: %s", err)
		}
	} else {
		log.Println("No env file detected.")
	}

}
