package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

/* HANDLER SUPPORT */
var CommandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
var ComponentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "subscribe",
		Description: "Subscribe to get notifications when Google Form opens",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "enabled",
				Description: "toggle for registration",
				Required:    true,
			},
		},
	},
}

func (runTimeData *RunTimeData) InitHandlers() {
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"subscribe": func(s *discordgo.Session, i *discordgo.InteractionCreate) { runTimeData.SubscribeHandler(s, i) },
	}
	ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"seen": func(s *discordgo.Session, i *discordgo.InteractionCreate) { runTimeData.SeenHandler(s, i)},
	}
}

/* EVENT HANDLERS */
func (runTimeData *RunTimeData) OnReady(s *discordgo.Session, r *discordgo.Ready) {
	scheduler.AddFunc(runTimeData.StartTime, runTimeData.startTicker)
	scheduler.AddFunc(runTimeData.StopTime, stopTicker)
	scheduler.Start()
	s.UpdateCustomStatus("👁️‍🗨️ Monitoring...")
	fmt.Println("***************************** BOT RUNNING *****************************")
}

func (runTimeData *RunTimeData) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if command, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			command(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if command, ok := ComponentHandlers[i.MessageComponentData().CustomID]; ok {
			command(s, i)
		}
	}
}

/* COMMAND HANDLERS */
func (runTimeData *RunTimeData) SubscribeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var memberID string
	if i.Interaction.GuildID != "" {
		memberID = i.Interaction.Member.User.ID
	} else {
		memberID = i.Interaction.User.ID
	}
	opt := i.ApplicationCommandData().Options[0].BoolValue()
	subscribed := runTimeData.GetSubscribed()
	if _, exists := subscribed[memberID]; opt && exists {
		SendMessageResponse(s, i, "You are already subscribed!")
	} else if opt {
		SendMessageResponse(s, i, "You have sucessfully subscribed!")
		fmt.Printf("SUCCESS @ %s : USER %s SUBSCRIBED\n", time.Now().Format("2006-01-02 15:04:05.00000"), memberID)
		runTimeData.InsertIntoDb(memberID)
	} else if _, exists := subscribed[memberID]; !opt && !exists {
		SendMessageResponse(s, i, "You are not currently subscribed!")
	} else if !opt {
		SendMessageResponse(s, i, "You have sucessfully unsubscribed!")
		runTimeData.RemoveFromDb(memberID)
	}
}

func (runTimeData *RunTimeData) SeenHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	memberID := i.Interaction.User.ID
	runTimeData.UpdateSeen(memberID, true)
	if runTimeData.GetSeen(memberID) {
		SendMessageResponse(s, i, "You have already marked this notification as seen.")
	} else {
		SendMessageResponse(s, i, "You have marked this notification as seen.")
	}
	fmt.Printf("SUCCESS @ %s : USER %s MARKED NOTIFICATION AS SEEN\n", time.Now().Format("2006-01-02 15:04:05.00000"), memberID)
}

/* MESSAGING */
func SendMessageResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
}

func SendComplexMessage(user *discordgo.User, dmChannel *discordgo.Channel, embed *discordgo.MessageEmbed, buttons ...discordgo.Button) error {
	components := make([]discordgo.MessageComponent, len(buttons))
	for i, button := range buttons {
		components[i] = button
	}
	_, err := s.ChannelMessageSendComplex(
		dmChannel.ID,
		&discordgo.MessageSend{
			Content: user.Mention(),
			Embeds:  []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	)
	return err
}

/* EMBEDS AND BUTTONS */
func OpenEmbed() *discordgo.MessageEmbed {
	now := time.Now()
	return &discordgo.MessageEmbed{
		Title: "Attendance Form Open!",
		Description: fmt.Sprintf("Today is `%s, %s %d, %d`\nIntro Data Science recitation attendance form has opened.\n",
			now.Weekday(), now.Month().String(), now.Day(), now.Year()),
		Color: 0x5865f2,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) FormLink() discordgo.Button {
	return discordgo.Button{
		Label: "Link",
		Style: discordgo.LinkButton,
		URL:   runTimeData.Url,
	}
}

func (runTimeData *RunTimeData) SeenButton() discordgo.Button {
	return discordgo.Button{
		Label: "Seen",
		Style: discordgo.PrimaryButton,
		CustomID: "seen",
	}
}
