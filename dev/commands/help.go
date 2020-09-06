package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
)

func Help(s disgord.Session, m *disgord.MessageCreate) {
	if !strings.HasPrefix(m.Message.Content, "l?help") || m.Message.Author.Bot {
		return
	}

	general := "avatar - [l?av]\nhelp - [l?help]\nserver info - [l?server]\nping test - [l?ping]"

	moderation := "ban - [l?ban]\nunban - [l?unban]\nkick - [l?kick]\nprune - [l?prune]\nlockdown - [l?lockdown]\nunlock - [l?unlock]"

	embed := &disgord.Embed{
		Color: 0xDDA1A1,
		Title: "Commands",
		Fields: []*disgord.EmbedField{
			{
				Name:  "General",
				Value: fmt.Sprintf("```ini\n%v\n```", general),
			},
			{
				Name:  "Moderation",
				Value: fmt.Sprintf("```ini\n%v\n```", moderation),
			},
		},
	}

	m.Message.Reply(context.Background(), s, embed)
}
