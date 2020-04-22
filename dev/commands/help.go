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

	content := "avatar - [l?av]\nhelp - [l?help]\nserver info - [l?server]\nping - [l?ping]"

	embed := &disgord.Embed{
		Color:       0xDDA1A1,
		Title:       "Commands",
		Description: fmt.Sprintf("```ini\n%v\n```", content),
	}

	m.Message.Reply(context.Background(), s, embed)
}
