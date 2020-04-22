package commands

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func HighestRole(s disgord.Session, gid, uid disgord.Snowflake) int {

	mem, err := s.GetMember(context.Background(), gid, uid)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	groles, err := s.GetGuildRoles(context.Background(), gid)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	sort.Sort(RoleByPos(groles))

	for _, gr := range groles {
		for _, r := range mem.Roles {
			if r == gr.ID {
				return gr.Position
			}
		}
	}

	return -1
}

func Ban(s disgord.Session, m *disgord.MessageCreate) {
	if !strings.HasPrefix(m.Message.Content, "l?ban") || m.Message.Author.Bot {
		return
	}

	cu, err := s.GetCurrentUser(context.Background())
	if err != nil {
		fmt.Println(err)
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
		return
	}

	args := getArgs(m.Message.Content)

	if len(args) < 1 {
		return
	}

	var (
		targetUser *disgord.User
		reason     string
		pruneDays  int
		//err        error
	)

	if len(args) == 1 {
		pruneDays = 0
		reason = ""
	} else if len(args) >= 2 {
		pruneDays, err = strconv.Atoi(args[1])
		if err != nil {
			pruneDays = 0
			reason = strings.Join(args[1:], " ")
		} else {
			reason = strings.Join(args[2:], " ")
		}
		if pruneDays > 7 {
			pruneDays = 7
		} else if pruneDays < 0 {
			pruneDays = 0
		}
	}

	if len(m.Message.Mentions) > 0 {
		targetUser = m.Message.Mentions[0]
	} else {
		sn, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return
		}
		targetUser, err = s.GetUser(context.Background(), disgord.NewSnowflake(sn))
		if err != nil {
			return
		}
	}

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
				userchannel.SendMsgString(context.Background(), s, fmt.Sprintf("You have been banned from %v", g.Name))

			} else {
				userchannel.SendMsgString(context.Background(), s, fmt.Sprintf("You have been banned from %v for the following reason:\n%v", g.Name, reason))
			}
		}
	}
	err = s.BanMember(context.Background(), m.Message.GuildID, targetUser.ID, &disgord.BanMemberParams{
		DeleteMessageDays: pruneDays,
		Reason:            fmt.Sprintf("%v: %v", m.Message.Author.Tag(), reason),
	})
	if err != nil {
		return
	}

	embed := &disgord.Embed{
		Title: "User banned",
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
