package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func Unban(s disgord.Session, m *disgord.MessageCreate) {
	if !strings.HasPrefix(m.Message.Content, "l?unban") || !strings.HasPrefix(m.Message.Content, "l?ub") || m.Message.Author.Bot {
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
		return
	}

	args := getArgs(m.Message.Content)

	if len(args) < 1 {
		return
	}

	userID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return
	}

	err = s.UnbanMember(context.Background(), m.Message.GuildID, disgord.NewSnowflake(userID), m.Message.Author.Tag())
	if err != nil {
		return
	}

	targetUser, err := s.GetUser(context.Background(), disgord.NewSnowflake(userID))
	if err != nil {
		return
	}

	embed := &disgord.Embed{
		Description: fmt.Sprintf("**Unbanned** %v - %v#%v (%v)", targetUser.Mention(), targetUser.Username, targetUser.Discriminator, targetUser.ID),
		Color:       0x00C800,
	}

	s.SendMsg(context.Background(), m.Message.ChannelID, embed)
}
