package cot

import (
	"fmt"
	"github.com/coreywagehoft/go-tak/pkg/cotproto"
	"github.com/google/uuid"
	"time"
)

// BasicMsg constructs a base CoT message with default fields.
func BasicMsg(typ, uid string, stale time.Duration) *cotproto.TakMessage {
	return &cotproto.TakMessage{
		CotEvent: &cotproto.CotEvent{
			Type:      typ,
			Access:    "",
			Qos:       "",
			Opex:      "",
			Uid:       uid,
			SendTime:  TimeToMillis(time.Now()),
			StartTime: TimeToMillis(time.Now()),
			StaleTime: TimeToMillis(time.Now().Add(stale)),
			How:       HowDefault,
			Lat:       0,
			Lon:       0,
			Hae:       NotNum,
			Ce:        NotNum,
			Le:        NotNum,
			Detail:    nil,
		},
	}
}

// MakePing returns a CoT ping message for the given UID.
func MakePing(uid string) *cotproto.TakMessage {
	return BasicMsg(TypePing, uid+"-ping", 10*time.Second)
}

// MakePong returns a CoT pong response message.
func MakePong() *cotproto.TakMessage {
	msg := BasicMsg(TypePong, "takPong", 20*time.Second)
	msg.CotEvent.How = HowResponse
	return msg
}

// MakeOfflineMsg creates an offline CoT message with a pp link.
func MakeOfflineMsg(uid, typ string) *cotproto.TakMessage {
	msg := BasicMsg(TypeOffline, uuid.New().String(), 3*time.Minute)
	msg.CotEvent.How = HowResponse
	xd := NewXMLDetails()
	xd.AddPpLink(uid, typ, "")
	msg.CotEvent.Detail = &cotproto.Detail{XmlDetail: xd.AsXMLString()}
	return msg
}

// MakeDpMsg builds a CoT SPI message with contact details.
func MakeDpMsg(uid, typ, name string, lat, lon float64) *cotproto.TakMessage {
	msg := BasicMsg(TypeDpSpi, uid+".SPI1", 20*time.Second)
	msg.CotEvent.How = HowEvent
	msg.CotEvent.Lat = lat
	msg.CotEvent.Lon = lon
	xd := NewXMLDetails()
	xd.AddPpLink(uid, typ, "")
	msg.CotEvent.Detail = &cotproto.Detail{
		XmlDetail: xd.AsXMLString(),
		Contact:   &cotproto.Contact{Callsign: name},
	}
	return msg
}

// TeamMarkerOpts contains options for constructing a team marker message.
type TeamMarkerOpts struct {
	UID       string
	Callsign  string
	GroupName string
	Role      string
	Lat, Lon  float64
	Stale     time.Duration
}

// MakeTeamMarker returns a CoT team position message, validating inputs.
func MakeTeamMarker(o TeamMarkerOpts) (*cotproto.TakMessage, error) {
	if o.Lat < -90 || o.Lat > 90 || o.Lon < -180 || o.Lon > 180 {
		return nil, fmt.Errorf("invalid coordinates: (%f, %f)", o.Lat, o.Lon)
	}
	if o.Stale <= 0 {
		return nil, fmt.Errorf("stale duration must be >0")
	}

	msg := BasicMsg(TypeTeam, o.UID, o.Stale)
	msg.CotEvent.How = HowDefault
	msg.CotEvent.Lat = o.Lat
	msg.CotEvent.Lon = o.Lon
	msg.CotEvent.Hae = 0
	msg.CotEvent.Ce = NotNum
	msg.CotEvent.Le = NotNum
	msg.CotEvent.Detail = &cotproto.Detail{
		Contact: &cotproto.Contact{
			Callsign: o.Callsign,
		},
		Group: &cotproto.Group{
			Name: o.GroupName,
			Role: o.Role,
		},
	}
	return msg, nil
}
