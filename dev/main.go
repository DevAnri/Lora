package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/devanri/Lora/dev/commands"

	"github.com/andersfylling/disgord"
)

type Config struct {
	Token            string   `json: "token"`
	OwoToken         string   `json: "owo_token"`
	ConnectionString string   `json: "connection_string"`
	DDMogChannels    []uint64 `json "dm_log_channels"`
	OwnerIds         []string `json: "owner_ids"`
}

var (
	config Config
)

var (
	dmChannels = []uint64{702741795922247751}
)

func ForwardDMs(s disgord.Session, m *disgord.MessageCreate) {

	if m.Message.Author.Bot {
		return
	}

	ch, err := s.GetChannel(context.Background(), m.Message.ChannelID)
	if err != nil {
		return
	}

	if ch.Type != disgord.ChannelTypeDM {
		return
	}

	embed := &disgord.Embed{
		Color:       0xDDA1A1,
		Title:       fmt.Sprintf("Message from %v", m.Message.Author.Tag()),
		Description: m.Message.Content,
		Footer:      &disgord.EmbedFooter{Text: m.Message.Author.ID.String()},
		Timestamp:   m.Message.Timestamp,
	}
	if len(m.Message.Attachments) > 0 {
		embed.Image = &disgord.EmbedImage{URL: m.Message.Attachments[0].URL}
	}

	for _, id := range dmChannels {
		s.SendMsg(context.Background(), disgord.NewSnowflake(id), embed)
	}
}

func main() {

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("Config file not found.\nPlease press enter.")
	}
	json.Unmarshal(file, &config)

	// setting up bot
	client := disgord.New(disgord.Config{
		BotToken: config.Token,
		Logger:   disgord.DefaultLogger(false),
	})

	defer client.StayConnectedUntilInterrupted(context.Background())

	//client.On(disgord.EvtMessageCreate, Ping)
	addHandlers(client)
}

func addHandlers(c *disgord.Client) {
	c.On(disgord.EvtMessageCreate, commands.Ping)
	c.On(disgord.EvtMessageCreate, commands.Avatar)
	c.On(disgord.EvtMessageCreate, commands.Server)
	c.On(disgord.EvtMessageCreate, commands.Help)
	c.On(disgord.EvtMessageCreate, commands.Ban)
	c.On(disgord.EvtMessageCreate, commands.Unban)
	c.On(disgord.EvtMessageCreate, commands.Kick)
	c.On(disgord.EvtMessageCreate, commands.Prune)
	c.On(disgord.EvtMessageCreate, ForwardDMs)
	c.On(disgord.EvtReady, Ready)
}

func Ready(s disgord.Session, r *disgord.Ready) {

	s.UpdateStatus(&disgord.UpdateStatusPayload{
		Game: &disgord.Activity{
			Type: disgord.ActivityTypeListening,
			Name: "all of you",
		},
	})
}
