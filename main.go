package main

import (
	"fmt"
	"log"
	"os"

	bot "github.com/devgrohl/GoBot/discord"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	discordTokenID := goDotEnvVariable("DISCORD_TOKEN_ID")
	fmt.Printf("Attempting to start bot with Token ID: %s\n", discordTokenID)

	bot.BotToken = discordTokenID
	bot.Run()
}
