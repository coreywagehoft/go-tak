package cot

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtoRW(t *testing.T) {
	msg := MakeDpMsg("testuid", "a-f-G", "test", 10, 20)

	b, err := MakeProtoPacket(msg)
	require.NoError(t, err)

	msg1, err := ReadProto(bufio.NewReader(bytes.NewReader(b)))
	require.NoError(t, err)

	assert.Equal(t, "testuid.SPI1", msg1.GetCotEvent().GetUid())
	assert.Equal(t, "b-m-p-s-p-i", msg1.GetCotEvent().GetType())
	assert.InDelta(t, 10., msg1.GetCotEvent().GetLat(), 0.0001)
	assert.InDelta(t, 20., msg1.GetCotEvent().GetLon(), 0.0001)
}

func TestMakeProtoStreamPacket(t *testing.T) {
	msg := MakeDpMsg("testuid", "a-f-G", "test", 10, 20)

	b, err := MakeProtoStreamPacket(msg)
	require.NoError(t, err)

	// Verify header format: magic byte + varint length
	assert.Equal(t, magic, b[0], "first byte should be magic byte")

	// Read the packet back
	msg1, err := ReadProtoStream(bufio.NewReader(bytes.NewReader(b)))
	require.NoError(t, err)

	assert.Equal(t, "testuid.SPI1", msg1.GetCotEvent().GetUid())
	assert.Equal(t, "b-m-p-s-p-i", msg1.GetCotEvent().GetType())
	assert.InDelta(t, 10., msg1.GetCotEvent().GetLat(), 0.0001)
	assert.InDelta(t, 20., msg1.GetCotEvent().GetLon(), 0.0001)
}

func TestMakeProtoMeshPacket(t *testing.T) {
	msg := MakeDpMsg("testuid", "a-f-G", "test", 35.5, -117.2)

	tests := []struct {
		name    string
		version uint64
	}{
		{"Version 0", ProtoVersion0},
		{"Version 1", ProtoVersion1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := MakeProtoMeshPacket(msg, tt.version)
			require.NoError(t, err)

			// Verify header format: magic byte + version varint + magic byte
			assert.Equal(t, magic, b[0], "first byte should be magic byte")

			// Read the packet back
			msg1, version, err := ReadProtoMesh(bufio.NewReader(bytes.NewReader(b)))
			require.NoError(t, err)

			assert.Equal(t, tt.version, version, "version should match")
			assert.Equal(t, "testuid.SPI1", msg1.GetCotEvent().GetUid())
			assert.Equal(t, "b-m-p-s-p-i", msg1.GetCotEvent().GetType())
			assert.InDelta(t, 35.5, msg1.GetCotEvent().GetLat(), 0.0001)
			assert.InDelta(t, -117.2, msg1.GetCotEvent().GetLon(), 0.0001)
		})
	}
}

func TestMakeProtoMeshPacketV1(t *testing.T) {
	msg := MakeDpMsg("testuid", "a-f-G", "test", 40.0, -75.0)

	b, err := MakeProtoMeshPacketV1(msg)
	require.NoError(t, err)

	// Read the packet back
	msg1, version, err := ReadProtoMesh(bufio.NewReader(bytes.NewReader(b)))
	require.NoError(t, err)

	assert.Equal(t, ProtoVersion1, version, "version should be 1")
	assert.Equal(t, "testuid.SPI1", msg1.GetCotEvent().GetUid())
	assert.Equal(t, "b-m-p-s-p-i", msg1.GetCotEvent().GetType())
	assert.InDelta(t, 40.0, msg1.GetCotEvent().GetLat(), 0.0001)
	assert.InDelta(t, -75.0, msg1.GetCotEvent().GetLon(), 0.0001)
}

func TestReadProtoStreamWithGarbage(t *testing.T) {
	msg := MakeDpMsg("testuid", "a-f-G", "test", 10, 20)

	b, err := MakeProtoStreamPacket(msg)
	require.NoError(t, err)

	// Prepend garbage data
	garbage := []byte{0x01, 0x02, 0x03, 0x04}
	data := append(garbage, b...)

	// Should skip garbage and find magic byte
	msg1, err := ReadProtoStream(bufio.NewReader(bytes.NewReader(data)))
	require.NoError(t, err)

	assert.Equal(t, "testuid.SPI1", msg1.GetCotEvent().GetUid())
}

func TestReadProtoMeshWithGarbage(t *testing.T) {
	msg := MakeDpMsg("testuid", "a-f-G", "test", 10, 20)

	b, err := MakeProtoMeshPacketV1(msg)
	require.NoError(t, err)

	// Prepend garbage data
	garbage := []byte{0x01, 0x02, 0x03, 0x04}
	data := append(garbage, b...)

	// Should skip garbage and find first magic byte
	msg1, version, err := ReadProtoMesh(bufio.NewReader(bytes.NewReader(data)))
	require.NoError(t, err)

	assert.Equal(t, ProtoVersion1, version)
	assert.Equal(t, "testuid.SPI1", msg1.GetCotEvent().GetUid())
}

func TestMultipleStreamMessages(t *testing.T) {
	msg1 := MakeDpMsg("uid1", "a-f-G", "test1", 10, 20)
	msg2 := MakeDpMsg("uid2", "a-f-G", "test2", 30, 40)
	msg3 := MakeDpMsg("uid3", "a-f-G", "test3", 50, 60)

	b1, err := MakeProtoStreamPacket(msg1)
	require.NoError(t, err)
	b2, err := MakeProtoStreamPacket(msg2)
	require.NoError(t, err)
	b3, err := MakeProtoStreamPacket(msg3)
	require.NoError(t, err)

	// Concatenate all messages
	data := append(append(b1, b2...), b3...)
	reader := bufio.NewReader(bytes.NewReader(data))

	// Read first message
	readMsg1, err := ReadProtoStream(reader)
	require.NoError(t, err)
	assert.Equal(t, "uid1.SPI1", readMsg1.GetCotEvent().GetUid())

	// Read second message
	readMsg2, err := ReadProtoStream(reader)
	require.NoError(t, err)
	assert.Equal(t, "uid2.SPI1", readMsg2.GetCotEvent().GetUid())

	// Read third message
	readMsg3, err := ReadProtoStream(reader)
	require.NoError(t, err)
	assert.Equal(t, "uid3.SPI1", readMsg3.GetCotEvent().GetUid())
}

func TestProtoVersionConstants(t *testing.T) {
	assert.Equal(t, uint64(0), ProtoVersion0, "ProtoVersion0 should be 0")
	assert.Equal(t, uint64(1), ProtoVersion1, "ProtoVersion1 should be 1")
}

func TestBackwardCompatibility(t *testing.T) {
	msg := MakeDpMsg("testuid", "a-f-G", "test", 10, 20)

	// Test that MakeProtoPacket produces same result as MakeProtoStreamPacket
	b1, err := MakeProtoPacket(msg)
	require.NoError(t, err)

	b2, err := MakeProtoStreamPacket(msg)
	require.NoError(t, err)

	assert.Equal(t, b1, b2, "MakeProtoPacket should produce same output as MakeProtoStreamPacket")

	// Test that ReadProto works same as ReadProtoStream
	msg1, err := ReadProto(bufio.NewReader(bytes.NewReader(b1)))
	require.NoError(t, err)

	msg2, err := ReadProtoStream(bufio.NewReader(bytes.NewReader(b2)))
	require.NoError(t, err)

	assert.Equal(t, msg1.GetCotEvent().GetUid(), msg2.GetCotEvent().GetUid())
}

func TestReadXMLEvent(t *testing.T) {
	// Create an XML event according to Traditional Protocol
	event := XMLBasicMsg("a-f-G-U-C", "test-uid-123", time.Minute)
	event.How = "h-g-i-g-o"
	event.Point.Lat = 35.5
	event.Point.Lon = -117.2
	event.AddCallsign("TestCallsign", "*:-1:stcp", true)

	// Marshal to XML with proper header
	xmlBytes, err := MakeXMLEvent(event)
	require.NoError(t, err)

	// Should start with XML header
	assert.True(t, bytes.HasPrefix(xmlBytes, []byte("<?xml")), "should start with XML header")

	// Read it back
	reader := bufio.NewReader(bytes.NewReader(xmlBytes))
	readEvent, err := ReadXMLEvent(reader)
	require.NoError(t, err)

	assert.Equal(t, event.UID, readEvent.UID)
	assert.Equal(t, event.Type, readEvent.Type)
	assert.InDelta(t, event.Point.Lat, readEvent.Point.Lat, 0.0001)
	assert.InDelta(t, event.Point.Lon, readEvent.Point.Lon, 0.0001)
}

func TestReadXMLEventMultiple(t *testing.T) {
	// Test reading multiple consecutive XML events from a stream
	event1 := XMLBasicMsg("a-f-G-U-C", "uid-1", time.Minute)
	event1.Point.Lat = 10.0
	event2 := XMLBasicMsg("a-f-G-U-C", "uid-2", time.Minute)
	event2.Point.Lat = 20.0
	event3 := XMLBasicMsg("a-f-G-U-C", "uid-3", time.Minute)
	event3.Point.Lat = 30.0

	xml1, err := MakeXMLEvent(event1)
	require.NoError(t, err)
	xml2, err := MakeXMLEvent(event2)
	require.NoError(t, err)
	xml3, err := MakeXMLEvent(event3)
	require.NoError(t, err)

	// Concatenate according to protocol: no arbitrary newlines between events
	stream := append(append(xml1, xml2...), xml3...)
	reader := bufio.NewReader(bytes.NewReader(stream))

	// Read all three events
	read1, err := ReadXMLEvent(reader)
	require.NoError(t, err)
	assert.Equal(t, "uid-1", read1.UID)
	assert.InDelta(t, 10.0, read1.Point.Lat, 0.0001)

	read2, err := ReadXMLEvent(reader)
	require.NoError(t, err)
	assert.Equal(t, "uid-2", read2.UID)
	assert.InDelta(t, 20.0, read2.Point.Lat, 0.0001)

	read3, err := ReadXMLEvent(reader)
	require.NoError(t, err)
	assert.Equal(t, "uid-3", read3.UID)
	assert.InDelta(t, 30.0, read3.Point.Lat, 0.0001)
}

func TestMakeXMLEvent(t *testing.T) {
	event := XMLBasicMsg("b-m-p-s-p-i", "marker-123", time.Second*30)
	event.How = "h-e"
	event.Point.Lat = 40.7128
	event.Point.Lon = -74.0060

	xmlBytes, err := MakeXMLEvent(event)
	require.NoError(t, err)

	// Verify it starts with XML header
	assert.True(t, bytes.Contains(xmlBytes, []byte("<?xml")))
	// Verify it contains the event
	assert.True(t, bytes.Contains(xmlBytes, []byte("<event")))
	assert.True(t, bytes.Contains(xmlBytes, []byte("</event>")))
	// Verify UID is in there
	assert.True(t, bytes.Contains(xmlBytes, []byte("marker-123")))
}
