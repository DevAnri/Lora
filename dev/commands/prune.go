package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func Prune(s disgord.Session, m *disgord.MessageCreate) {
	if !strings.HasPrefix(m.Message.Content, "l?prune") || m.Message.Author.Bot {
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

	if botperms&disgord.PermissionManageMessages == 0 && botperms&disgord.PermissionAdministrator == 0 {
		return
	}

	uperms, err := s.GetMemberPermissions(context.Background(), m.Message.GuildID, m.Message.Author.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if uperms&disgord.PermissionManageMessages == 0 && uperms&disgord.PermissionAdministrator == 0 {
		return
	}

	args := getArgs(m.Message.Content)

	if len(args) < 1 {
		return
	}

	amt, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	if amt > 100 || amt < 1 {
		return
	}

	if amt < 2 {
		//delete just one message

		msgs, err := s.GetMessages(context.Background(), m.Message.ChannelID, &disgord.GetMessagesParams{
			Before: m.Message.ID,
			Limit:  1,
		})
		if err != nil {
			return
		}
		msgids := []disgord.Snowflake{m.Message.ID, msgs[0].ID}

		s.DeleteMessages(context.Background(), m.Message.ChannelID, &disgord.DeleteMessagesParams{
			Messages: msgids,
		})

	} else {
		//do it normally

		msgs, err := s.GetMessages(context.Background(), m.Message.ChannelID, &disgord.GetMessagesParams{
			Before: m.Message.ID,
			Limit:  uint(amt),
		})
		if err != nil {
			return
		}

		var msgids []disgord.Snowflake

		for _, msg := range msgs {
			msgids = append(msgids, msg.ID)
		}

		s.DeleteMessages(context.Background(), m.Message.ChannelID, &disgord.DeleteMessagesParams{
			Messages: msgids,
		})

		s.DeleteMessage(context.Background(), m.Message.ChannelID, m.Message.ID)

	}

	cmsg, err := s.SendMsg(context.Background(), m.Message.ChannelID, fmt.Sprintf("Removed %v message(s)", amt))
	if err != nil {
		return
	}

	<-time.After(time.Second * 3)

	s.DeleteMessage(context.Background(), cmsg.ChannelID, cmsg.ID)
}
