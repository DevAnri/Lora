package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func Kick(s disgord.Session, m *disgord.MessageCreate) {
	if !strings.HasPrefix(m.Message.Content, "l?kick") || !strings.HasPrefix(m.Message.Content, "l?k") || m.Message.Author.Bot {
		return
	}
	cu, err := s.GetCurrentUser(context.Background())
	if err != nil {
		return
	}

	botperms, err := s.GetMemberPermissions(context.Background(), m.Message.GuildID, cu.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if botperms&disgord.PermissionBanMembers == 0 && botperms&disgord.PermissionAdministrator == 0 {
		return
	}

	uperms, err := s.GetMemberPermissions(context.Background(), m.Message.GuildID, m.Message.Author.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if uperms&disgord.PermissionBanMembers == 0 && uperms&disgord.PermissionAdministrator == 0 {
		s.SendMsg(context.Background(), m.Message.ChannelID, "no")
		return
	}

	args := getArgs(m.Message.Content)

	if len(args) < 1 {
		return
	}

	var (
		targetUser *disgord.User
		reason     string
		//err        error
	)

	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}

	var userID disgord.Snowflake
	if len(m.Message.Mentions) > 0 {
		userID = m.Message.Mentions[0].ID
	} else {
		sn, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return
		}
		userID = disgord.NewSnowflake(sn)
	}

	u, err := s.GetMember(context.Background(), m.Message.GuildID, userID)
	if err != nil {
		// this means that the user is not in the server
		return
	}

	targetUser = u.User

	if targetUser.ID == cu.ID {
		return
	}
	if targetUser.ID == m.Message.Author.ID {
		return
	}

	topUserrole := HighestRole(s, m.Message.GuildID, m.Message.Author.ID)
	topTargetrole := HighestRole(s, m.Message.GuildID, targetUser.ID)
	topBotrole := HighestRole(s, m.Message.GuildID, cu.ID)

	if topUserrole <= topTargetrole || topBotrole <= topTargetrole {
		return
	}

	if topTargetrole > 0 {

		okCh := true

		userchannel, err := s.CreateDM(context.Background(), targetUser.ID)
		if err != nil {
			okCh = false
		}

		if okCh {
			g, err := s.GetGuild(context.Background(), m.Message.GuildID)
			if err != nil {
				return
			}

			if reason == "" {
				userchannel.SendMsgString(context.Background(), s, fmt.Sprintf("You have been kicked from %v", g.Name))

			} else {
				userchannel.SendMsgString(context.Background(), s, fmt.Sprintf("You have been kicked from %v for the following reason:\n%v", g.Name, reason))
			}
		}
	}
	err = s.KickMember(context.Background(), m.Message.GuildID, targetUser.ID, fmt.Sprintf("%v: %v", m.Message.Author.Tag(), reason))
	if err != nil {
		return
	}

	embed := &disgord.Embed{
		Title: "User kicked",
		Color: 0xC80000,
		Fields: []*disgord.EmbedField{
			{
				Name:   "Username",
				Value:  fmt.Sprintf("%v", targetUser.Mention()),
				Inline: true,
			},
			{
				Name:   "ID",
				Value:  fmt.Sprintf("%v", targetUser.ID),
				Inline: true,
			},
		},
	}

	s.SendMsg(context.Background(), m.Message.ChannelID, embed)

}
