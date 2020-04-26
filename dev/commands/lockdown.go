package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
)

func Lockdown(s disgord.Session, m *disgord.MessageCreate) {

	if !strings.HasPrefix(m.Message.Content, "l?lockdown") || m.Message.Author.Bot {
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

	if botperms&disgord.PermissionManageRoles == 0 && botperms&disgord.PermissionAdministrator == 0 {
		return
	}

	uperms, err := s.GetMemberPermissions(context.Background(), m.Message.GuildID, m.Message.Author.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if uperms&disgord.PermissionManageRoles == 0 && uperms&disgord.PermissionAdministrator == 0 {
		return
	}

	var (
		er *disgord.Role
		ep disgord.PermissionOverwrite
	)

	grs, err := s.GetGuildRoles(context.Background(), m.Message.GuildID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, gr := range grs {
		if gr.ID == m.Message.GuildID {
			er = gr
		}
	}

	ch, err := s.GetChannel(context.Background(), m.Message.ChannelID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ov := range ch.PermissionOverwrites {
		if ov.ID == er.ID {
			ep = ov
		}
	}

	if er == nil {
		return
	}

	if ep.ID.IsZero() {
		ep = disgord.PermissionOverwrite{
			Type:  "role",
			Allow: 0,
			Deny:  0,
			ID:    er.ID,
		}
	}

	if ep.Allow&disgord.PermissionSendMessages == 0 && ep.Deny&disgord.PermissionSendMessages == 0 {
		// DEFAULT
		err := s.UpdateChannelPermissions(
			context.Background(),
			m.Message.ChannelID,
			er.ID,
			&disgord.UpdateChannelPermissionsParams{
				Type:  "role",
				Allow: ep.Allow,
				Deny:  ep.Deny + disgord.PermissionSendMessages,
			},
		)
		if err != nil {
			s.SendMsg(context.Background(), m.Message.ChannelID, "Could not lock channel")
			return
		}
		s.SendMsg(context.Background(), m.Message.ChannelID, "Channel locked")
	} else if ep.Allow&disgord.PermissionSendMessages != 0 && ep.Deny&disgord.PermissionSendMessages == 0 {
		// IS ALLOWED
		err := s.UpdateChannelPermissions(
			context.Background(),
			m.Message.ChannelID,
			er.ID,
			&disgord.UpdateChannelPermissionsParams{
				Type:  "role",
				Allow: ep.Allow - disgord.PermissionSendMessages,
				Deny:  ep.Deny + disgord.PermissionSendMessages,
			},
		)
		if err != nil {
			s.SendMsg(context.Background(), m.Message.ChannelID, "Could not lock channel")
			return
		}
		s.SendMsg(context.Background(), m.Message.ChannelID, "Channel locked")
	} else if ep.Allow&disgord.PermissionSendMessages == 0 && ep.Deny&disgord.PermissionSendMessages != 0 {
		// IS DENIED
		s.SendMsg(context.Background(), m.Message.ChannelID, "Channel already locked")
	}
}
