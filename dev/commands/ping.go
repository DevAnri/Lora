package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func Ping(s disgord.Session, m *disgord.MessageCreate) {

	if !strings.HasPrefix(m.Message.Content, "l?ping") {
		return
	}
	st := time.Now()

	msg, err := m.Message.Reply(context.Background(), s, "pong")
	if err != nil {
		return
	}

	p := time.Now().Sub(st)
	s.UpdateMessage(context.Background(), msg.ChannelID, msg.ID).SetContent("pong\n" + fmt.Sprint(p)).Execute()
}
