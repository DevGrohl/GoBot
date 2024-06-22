package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID")
	BotToken       = flag.String("token", goDotEnvVariable("DISCORD_TOKEN_ID"), "Bot access Token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands aftger shutdown")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	checkNilErr(err, "New Session")
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func checkNilErr(e error, msg ...string) {
	if e != nil {
		log.Fatal(strings.Join(msg, " ") + " Fatal error: " + e.Error())
	}
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	checkNilErr(err, "Error opening connection")

	moderator := moderator(s, *GuildID)

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Shutting down...")
}
