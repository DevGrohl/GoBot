package bot // github.com/devgrohl/GoBot/internal/discord/bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {
	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)

	// open Session
	discord.Open()
	defer discord.Close() // close session after func termination

	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func showInfo(discord *discordgo.Session) {
	// Print bot information
	fmt.Println("Bot Username: " + discord.State.User.Username)
	fmt.Println("Bot ID: " + discord.State.User.ID)

	// Print server information
	for _, s := range discord.State.Guilds {
		fmt.Println("------------- Server Information ---------------")
		fmt.Println("Server Count: " + strconv.Itoa(len(discord.State.Guilds)))
		fmt.Println("Server Name: " + s.Name)
		fmt.Println("Server ID: " + s.ID)
		fmt.Println("Server Owner: " + s.OwnerID)
		fmt.Println("Total members in server: " + strconv.Itoa(s.MemberCount))
		fmt.Println("Total roles in server: " + strconv.Itoa(len(s.Roles)))
	}
}

func timeoutUser(discord *discordgo.Session, channelID string, userID string) {
	discord.ChannelMessageSend(channelID, "Timeout user: "+userID+"Channel ID: "+channelID)

	// Check for permissions to Timeout
	perms, err := discord.UserChannelPermissions(discord.State.User.ID, channelID)
	checkNilErr(err)
	if perms&discordgo.PermissionAdministrator != 0 {
		discord.ChannelMessageSend(channelID, "I don't have permissions to timeout this user")
	}

	// Timeout user for 60 seconds
	t := time.Now().Add(time.Second * 60)
	err = discord.GuildMemberTimeout(discord.State.Guilds[0].ID, userID, &t)
	fmt.Println("Timing out user: " + userID + "for: " + t.String() + " seconds")
	if err != nil {
		discord.ChannelMessageSend(channelID, "Error: "+err.Error())
	}
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World!")
	case strings.Contains(message.Content, "!watch"):
		discord.ChannelMessageSend(message.ChannelID, "I'm watching !")
	case strings.Contains(message.Content, "!showInfo"):
		showInfo(discord)
	case strings.Contains(message.Content, "!timeout"):
		// get mentioned user
		user := message.Mentions[0]
		timeoutUser(discord, message.ChannelID, user.ID)
	case strings.Contains(message.Content, "!roles"):
		// get all roles
		roles := discord.State.Guilds[0].Roles
		for _, r := range roles {
			discord.ChannelMessageSend(message.ChannelID, r.Name)
		}
	case strings.Contains(message.Content, "!rm_role"):
		// rm role from mentioned User
		user := message.Mentions[0]
		role := message.MentionRoles
		discord.ChannelMessageSend(message.ChannelID, "Removing role: "+role[0]+" from user: "+user.Username)
		err := discord.GuildMemberRoleRemove(discord.State.Guilds[0].ID, user.ID, role[0])
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Error: "+err.Error())
		}
	}
}
