package squadrcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// PacketHeaderSize is the Size of the header in bytes excluding the Size field.
	PacketHeaderSize int32 = 8

	// PacketAmountOfNullTerminators is the amount of null terminators at the end of the packet.
	PacketAmountOfNullTerminators int32 = 2
	maxPacketSize                       = 4096
)

type packet struct {

	// Size is a 32-bit little endian integer, representing the length of the request in bytes.
	// Note that the packet Size field itself is not included when determining the Size of the
	// packet, so the value of this field is always 4 less than the packet's actual length.
	//
	Size int32

	Id int32

	Type int32

	Body []byte
}

func newPacket(packetType int32, id int32, body string) *packet {
	return &packet{
		Size: int32(len([]byte(body))) + PacketHeaderSize + PacketAmountOfNullTerminators,
		Type: packetType,
		Id:   id,
		Body: []byte(body),
	}
}

// GetBodySizeWithTerminators returns the size of the body of the packet in bytes.
func (packet *packet) GetBodySizeWithTerminators() int32 {
	return packet.Size - 4 - 4 // minus sizeof(Id) and sizeof(Type)
}

// GetBodySize returns the size of the body of the packet in bytes without padding.
func (packet *packet) GetBodySize() int32 {
	return packet.GetBodySizeWithTerminators() - PacketAmountOfNullTerminators
}

func (packet *packet) GetBody() string {
	return string(packet.Body)
}

func (packet *packet) ReadFrom(r io.Reader) (int64, error) {
	reader := &countingReader{
		Reader: r,
	}
	r = reader // Prevent non-counting reader from being used accidentally in rest of function

	if err := binary.Read(reader, binary.LittleEndian, &packet.Size); err != nil {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("failure to read packet Size: %w", err),
			PacketBytes: reader.Bytes,
		}
	}

	fmt.Printf("Packet size %v\n", packet.Size)

	if packet.Size > maxPacketSize {
		panic(fmt.Errorf("packet size too large, %d", packet.Size))
		//return reader.TotalBytesRead, &packetParseError{
		//	Err:         fmt.Errorf("packet size too large, %d", packet.Size),
		//	PacketBytes: reader.Bytes,
		//}
	}

	if err := binary.Read(reader, binary.LittleEndian, &packet.Id); err != nil {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("failure to read packet id: %w", err),
			PacketBytes: reader.Bytes,
		}
	}

	if err := binary.Read(reader, binary.LittleEndian, &packet.Type); err != nil {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("failure to read packet type: %w", err),
			PacketBytes: reader.Bytes,
		}
	}

	bodySize := packet.GetBodySize()

	if bodySize < 0 {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("body Size < 0: %d", bodySize),
			PacketBytes: reader.Bytes,
		}
	}

	packet.Body = make([]byte, bodySize)
	_, err := reader.Read(packet.Body)

	if err != nil {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("error parsing body, body Size: %d. %w", bodySize, err),
			PacketBytes: reader.Bytes,
		}
	}

	bodyTerminator := make([]byte, 1)
	_, err = reader.Read(bodyTerminator)

	if err != nil {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("error reading body nul terminator. %w", err),
			PacketBytes: reader.Bytes,
		}
	}

	if bodyTerminator[0] != 0 {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("body terminator is not nul, %x", bodyTerminator[0]),
			PacketBytes: reader.Bytes,
		}
	}

	packetTerminator := make([]byte, 1)
	_, err = reader.Read(packetTerminator)

	if err != nil {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("error reading packet terminator. %w", err),
			PacketBytes: reader.Bytes,
		}
	}

	if packetTerminator[0] != 0 {
		return reader.TotalBytesRead, &packetParseError{
			Err:         fmt.Errorf("packet terminator is not nul, %x", packetTerminator[0]),
			PacketBytes: reader.Bytes,
		}
	}

	fmt.Printf("Packet read: %v\n", reader.Bytes)

	return reader.TotalBytesRead, nil
}

func (packet *packet) WriteTo(w io.Writer) (int64, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, packet.Size+4))

	if err := binary.Write(buffer, binary.LittleEndian, packet.Size); err != nil {
		return 0, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, packet.Id); err != nil {
		return 0, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, packet.Type); err != nil {
		return 0, err
	}

	buffer.Write(packet.Body)

	for i := int32(0); i < PacketAmountOfNullTerminators; i++ {
		buffer.WriteByte(0x00)
	}

	return buffer.WriteTo(w)
}

type packetParseError struct {
	Err         error
	PacketBytes []byte
}

func (e *packetParseError) Error() string {
	return fmt.Sprintf("failed to parse packet: %s. Bytes: % x", e.Err, e.PacketBytes)
}

func (e *packetParseError) Unwrap() error {
	return e.Err
}
