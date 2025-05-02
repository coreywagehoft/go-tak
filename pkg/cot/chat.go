package cot

import (
	"fmt"
	"html"
	"time"

	"github.com/coreywagehoft/go-tak/pkg/cotproto"
)

type ChatMessage struct {
	ID       string    `json:"message_id"`
	Time     time.Time `json:"time"`
	Parent   string    `json:"parent"`
	Chatroom string    `json:"chatroom"`
	From     string    `json:"from"`
	FromUID  string    `json:"from_uid"`
	ToUID    string    `json:"to_uid"`
	Direct   bool      `json:"direct"`
	Text     string    `json:"text"`
}

func MakeChatMessage(c *ChatMessage) *cotproto.TakMessage {
	t := time.Now().UTC().Format(time.RFC3339)
	msgUID := fmt.Sprintf("GeoChat.%s.%s.%s", c.FromUID, c.ToUID, c.ID)
	msg := BasicMsg("b-t-f", msgUID, time.Second*10)
	msg.CotEvent.How = "h-g-i-g-o"
	xd := NewXMLDetails()
	xd.AddPpLink(c.FromUID, "", "")

	chat := xd.AddOrChangeChild("__chat", map[string]string{"parent": c.Parent, "groupOwner": "false", "chatroom": c.Chatroom, "senderCallsign": c.From, "id": c.ToUID, "messageId": c.ID})
	chat.AddOrChangeChild("chatgrp", map[string]string{"uid0": c.FromUID, "uid1": c.ToUID, "id": c.ToUID})

	xd.AddChild("remarks", map[string]string{"source": "BAO.F.ATAK." + c.FromUID, "to": c.ToUID, "time": t}, html.EscapeString(c.Text))

	if c.Direct {
		marti := xd.AddChild("marti", nil, "")
		marti.AddChild("dest", map[string]string{"callsign": c.Chatroom}, "")
	}

	msg.CotEvent.Detail = &cotproto.Detail{XmlDetail: xd.AsXMLString()}

	return msg
}

func MsgToChat(m *CotMessage) *ChatMessage {
	chat := m.GetDetail().GetFirst("__chat")
	if chat == nil {
		return nil
	}

	c := &ChatMessage{
		ID:       chat.GetAttr("messageId"),
		Time:     m.GetStartTime(),
		Parent:   chat.GetAttr("parent"),
		Chatroom: chat.GetAttr("chatroom"),
		From:     chat.GetAttr("senderCallsign"),
		ToUID:    chat.GetAttr("id"),
	}

	if cg := chat.GetFirst("chatgrp"); cg != nil {
		c.FromUID = cg.GetAttr("uid0")
	}

	if link := m.GetFirstLink("p-p"); link != nil {
		if uid := link.GetAttr("uid"); uid != "" {
			c.FromUID = uid
		}
	}

	if c.Chatroom != c.ToUID {
		c.Direct = true
	}

	if rem := m.GetDetail().GetFirst("remarks"); rem != nil {
		c.Text = html.UnescapeString(rem.GetText())
	} else {
		return nil
	}

	return c
}