package commands

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func getArgs(s string) []string {
	f := strings.Fields(s)
	if len(f) > 1 {
		return f[1:]
	}
	return []string{}
}
func Avatar(s disgord.Session, m *disgord.MessageCreate) {

	args := getArgs(m.Message.Content)

	if !strings.HasPrefix(m.Message.Content, "l?av") || m.Message.Author.Bot {
		return
	}

	var (
		targetUser *disgord.User
	)
	if len(args) > 0 {
		if len(m.Message.Mentions) > 0 {
			targetUser = m.Message.Mentions[0]
		} else {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return
			}
			targetUser, err = s.GetUser(context.Background(), disgord.NewSnowflake(id))
			if err != nil {
				return
			}
		}
	} else {
		targetUser = m.Message.Author
	}

	url, _ := targetUser.AvatarURL(2048, true)

	embed := &disgord.Embed{
		Color: HighestColor(s, m.Message.GuildID, targetUser.ID),
		Title: targetUser.Tag(),
		Image: &disgord.EmbedImage{
			URL: url,
		},
	}
	m.Message.Reply(context.Background(), s, embed)
}

func HighestColor(s disgord.Session, gid, uid disgord.Snowflake) int {

	mem, err := s.GetMember(context.Background(), gid, uid)
	if err != nil {
		return 0
	}

	groles, err := s.GetGuildRoles(context.Background(), gid)
	if err != nil {
		return 0
	}

	sort.Sort(RoleByPos(groles))

	for _, gr := range groles {
		for _, r := range mem.Roles {
			if r == gr.ID {
				if gr.Color != 0 {
					return int(gr.Color)
				}
			}
		}
	}

	return 0
}

type RoleByPos []*disgord.Role

func (a RoleByPos) Len() int           { return len(a) }
func (a RoleByPos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RoleByPos) Less(i, j int) bool { return a[i].Position > a[j].Position }
