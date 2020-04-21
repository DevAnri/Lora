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

	g, err := s.GetGuild(context.Background(), m.Message.GuildID)
	if err != nil {
		fmt.Println(err)
		return
	}
	u, err := s.GetCurrentUser(context.Background())
	if err != nil {

	}
	pfp, err := u.AvatarURL(2048, true)
	if err != nil {
		return
	}

	embed := &disgord.Embed{
		Color:       0xDDA1A1,
		Fields: []*disgord.EmbedField{

		},
	}
	if g.Icon != "" {
		embed.Author = &disgord.EmbedAuthor{
			Name:    "Help",
			IconURL: pfp,
		}
	}
	m.Message.Reply(context.Background(), s, embed)
}
