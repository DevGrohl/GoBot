package moderation

import (
	"fmt"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func main(s *discordgo.Session, GuildID string) {
	s.Identify.Intents |= discordgo.IntentAutoModerationExecution
	s.Identify.Intents |= discordgo.IntentMessageContent

	enabled := true
	rule, err := s.AutoModerationRuleCreate(GuildID, &discordgo.AutoModerationRule{
		Name:        "Auto Moderation",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"negro"},
		},

		Enabled: &enabled,
		Actions: []discordgo.AutoModerationAction{
			{
				Type: discordgo.AutoModerationRuleActionBlockMessage,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Auto moderation rules created")
	defer s.AutoModerationRuleDelete(GuildID, rule.ID)

	s.AddHandlerOnce(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
		_, err = s.AutoModerationRuleEdit(GuildID, rule.ID, &discordgo.AutoModerationRule{
			TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
				KeywordFilter: []string{"negro"},
			},

			Actions: []discordgo.AutoModerationAction{
				{Type: discordgo.AutoModerationRuleActionTimeout, Metadata: &discordgo.AutoModerationActionMetadata{Duration: 60}},
				{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{ChannelID: e.ChannelID}},
			},
		})
		if err != nil {
			s.AutoModerationRuleDelete(GuildID, rule.ID)
			panic(err)
		}

		var counter int
		var counterMutex sync.Mutex

		s.AddHandler(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
			action := "UNK"

			switch e.Action.Type {
			case discordgo.AutoModerationRuleActionBlockMessage:
				action = "Block Message"
			case discordgo.AutoModerationRuleActionSendAlertMessage:
				action = "Send Alert Message"
			case discordgo.AutoModerationRuleActionTimeout:
				action = "Timeout User"
			}

			counterMutex.Lock()
			counter++

			if counter == 1 {
				counterMutex.Unlock()
				s.ChannelMessageSend(e.ChannelID, "First occurrence")
			} else if counter == 2 {
				counterMutex.Unlock()
				s.ChannelMessageSend(e.ChannelID, "Multiple occurrences of the same violation")

				s.Close()
				s.AutoModerationRuleDelete(GuildID, rule.ID)
				os.Exit(0)
			}
		})
	})
}
