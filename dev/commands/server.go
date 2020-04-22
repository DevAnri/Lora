package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
)

func Server(s disgord.Session, m *disgord.MessageCreate) {

	if !strings.HasPrefix(m.Message.Content, "l?server") || m.Message.Author.Bot {
		return
	}
	g, err := s.GetGuild(context.Background(), m.Message.GuildID)
	if err != nil {
		return
	}

	embed := &disgord.Embed{
		Color: 0xDDA1A1,
		Fields: []*disgord.EmbedField{
			{
				Name:   "Owner",
				Value:  fmt.Sprintf("<@!%v>", g.OwnerID),
				Inline: true,
			},
			{
				Name:  "Server ID",
				Value: fmt.Sprintln(g.ID),
			},
			{
				Name:   "Members",
				Value:  fmt.Sprintln(g.MemberCount),
				Inline: true,
			},
			{
				Name:   "Channels",
				Value:  fmt.Sprintln(len(g.Channels)),
				Inline: true,
			},
			{
				Name:   "Emotes",
				Value:  fmt.Sprintln(len(g.Emojis)),
				Inline: true,
			},
		},
	}

	if g.Icon != "" {
		embed.Thumbnail = &disgord.EmbedThumbnail{
			URL: fmt.Sprintf("https://cdn.discordapp.com/icons/%v/%v.png", g.ID, g.Icon),
		}
		embed.Author = &disgord.EmbedAuthor{
			Name:    g.Name,
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/icons/%v/%v.png", g.ID, g.Icon),
		}
	}
	m.Message.Reply(context.Background(), s, embed)
}
