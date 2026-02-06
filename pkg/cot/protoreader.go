package cot

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/coreywagehoft/go-tak/pkg/cotproto"
)

const (
	magic byte = 0xbf

	// TAK Protocol Versions
	ProtoVersion0 uint64 = 0 // Traditional XML CoT
	ProtoVersion1 uint64 = 1 // Protobuf-based CoT
)

type ProtoReader struct {
	r *bufio.Reader
}

func NewProtoReader(r io.Reader) *ProtoReader {
	if rb, ok := r.(*bufio.Reader); ok {
		return &ProtoReader{r: rb}
	}

	return &ProtoReader{r: bufio.NewReader(r)}
}

func (er *ProtoReader) ReadProtoBuf() (*cotproto.TakMessage, error) {
	return ReadProtoStream(er.r)
}

// MakeProtoStreamPacket creates a TAK Protocol Streaming message for streaming connections.
// Format: <magic byte 0xbf> <message length varint> <payload>
func MakeProtoStreamPacket(msg *cotproto.TakMessage) ([]byte, error) {
	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len(payload)+10) // max varint is 10 bytes
	buf[0] = magic
	n := binary.PutUvarint(buf[1:], uint64(len(payload)))
	copy(buf[1+n:], payload)

	return buf[:1+n+len(payload)], nil
}

// MakeProtoMeshPacket creates a TAK Protocol message for mesh networks (UDP/directed TCP).
// Format: <magic byte 0xbf> <protocol version varint> <magic byte 0xbf> <payload>
func MakeProtoMeshPacket(msg *cotproto.TakMessage, version uint64) ([]byte, error) {
	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len(payload)+12) // 2 magic bytes + max varint (10 bytes)
	buf[0] = magic
	n := binary.PutUvarint(buf[1:], version)
	buf[1+n] = magic
	copy(buf[2+n:], payload)

	return buf[:2+n+len(payload)], nil
}

// MakeProtoMeshPacketV1 creates a TAK Protocol Version 1 message for mesh networks.
// This is a convenience function that uses ProtoVersion1.
func MakeProtoMeshPacketV1(msg *cotproto.TakMessage) ([]byte, error) {
	return MakeProtoMeshPacket(msg, ProtoVersion1)
}

// MakeProtoPacket is deprecated. Use MakeProtoStreamPacket for streaming connections
// or MakeProtoMeshPacket for mesh networks.
func MakeProtoPacket(msg *cotproto.TakMessage) ([]byte, error) {
	return MakeProtoStreamPacket(msg)
}

func ReadProtoStream(r *bufio.Reader) (*cotproto.TakMessage, error) {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == magic {
			break
		}
	}

	size, err := binary.ReadUvarint(r)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(r, buf)

	if err != nil {
		return nil, err
	}

	msg := new(cotproto.TakMessage)
	err = proto.Unmarshal(buf, msg)

	return msg, err
}

// ReadProtoMesh reads a TAK Protocol message for mesh networks.
// Format: <magic byte 0xbf> <protocol version varint> <magic byte 0xbf> <payload>
// Note: For UDP, reader should contain only a single datagram.
// For directed TCP messages (connect, send, disconnect), the connection should
// close after the message, making the rest of the stream the payload.
func ReadProtoMesh(r *bufio.Reader) (*cotproto.TakMessage, uint64, error) {
	// Read first magic byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		if b == magic {
			break
		}
	}

	// Read protocol version
	version, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, 0, err
	}

	// Read second magic byte
	b, err := r.ReadByte()
	if err != nil {
		return nil, 0, err
	}

	if b != magic {
		return nil, 0, io.ErrUnexpectedEOF
	}

	// For mesh messages, the spec doesn't define a length field.
	// The payload is the remainder of the message:
	// - For UDP: rest of the datagram
	// - For directed TCP: rest of data until connection closes
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, 0, err
	}

	msg := new(cotproto.TakMessage)
	err = proto.Unmarshal(buf, msg)

	return msg, version, err
}

// ReadProto is deprecated. Use ReadProtoStream for streaming connections
// or ReadProtoMesh for mesh network messages.
func ReadProto(r *bufio.Reader) (*cotproto.TakMessage, error) {
	return ReadProtoStream(r)
}

// ReadXMLEvent reads a Traditional Protocol (Version 0) XML CoT event from a stream.
// According to the spec, XML events are delimited by searching for "</event>" token.
// Messages are prefaced by XML header (<?xml ... ?>) followed by newline.
func ReadXMLEvent(r *bufio.Reader) (*Event, error) {
	var buf bytes.Buffer
	var inEvent bool
	var foundStart bool

	// Read until we find and consume the complete </event> tag
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		buf.WriteByte(b)

		// Look for start of event tag
		if !foundStart && buf.Len() >= 6 {
			recent := buf.Bytes()
			// Check for "<event" anywhere in the buffer
			if bytes.Contains(recent, []byte("<event")) {
				foundStart = true
				inEvent = true
			}
		}

		// Look for closing tag
		if inEvent && b == '>' && buf.Len() >= 8 {
			recent := buf.Bytes()
			if bytes.HasSuffix(recent, []byte("</event>")) {
				break
			}
		}

		if buf.Len() > 1024*1024 { // 1MB safety limit
			return nil, fmt.Errorf("message too large")
		}
	}

	// Parse the XML - extract just the event part
	eventStart := bytes.Index(buf.Bytes(), []byte("<event"))
	if eventStart == -1 {
		return nil, fmt.Errorf("no event tag found")
	}

	eventXML := buf.Bytes()[eventStart:]
	var event Event
	err := xml.Unmarshal(eventXML, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML event: %w", err)
	}

	return &event, nil
}
