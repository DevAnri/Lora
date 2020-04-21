package main

import (
	"Lora/dev/commands"
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/andersfylling/disgord"
)

type Config struct {
	Token string `json:"token"`
}

func main() {

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("Config file not found.\nPlease press enter.")
	}
	var config Config
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
}
