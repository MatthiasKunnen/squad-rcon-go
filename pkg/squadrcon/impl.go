package squadrcon

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Packet type
const (
	serverDataResponseValue = 0
	serverDataExecCommand   = 2
	serverDataAuthResponse  = 2
	serverDataAuth          = 3
)

const (
	wrapIdsAfter = 200
)

var (
	ErrCommandEmpty = errors.New("command is empty")
)

type rconResponse struct {
	Body []byte
}

type callback struct {
	// The channel which will be passed the data as responded by RCON.
	Channel chan string

	// The aggregated data.
	Data []byte

	// Whether the response could be spread out over multiple packets.
	// When true, a confirmation command should be sent.
	MaybeMultiPacket bool
}

type rconImpl struct {
	authenticated bool
	callbackLock  sync.Mutex
	callbacks     map[int32]*callback
	conn          net.Conn
	dialTimeout   time.Duration
	execIdCounter int
	idCounterLock sync.Mutex
	startId       int
	writeTimeout  time.Duration
}

// Close closes the connection.
func (r *rconImpl) Close() error {
	if err := r.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (r *rconImpl) Execute(command string) (string, error) {
	if command == "" {
		return "", ErrCommandEmpty
	}

	packetId := r.getNextId()
	fmt.Printf("Executing command %s, ID: %d\n", command, packetId)

	if err := r.write(serverDataExecCommand, packetId, command); err != nil {
		return "", err
	}

	// Send a short command with packetId + 1 after each command. When we receive the response to
	// this command, we know that all the responses of the previous command have arrived.
	if err := r.write(serverDataExecCommand, packetId+1, "ShowCurrentMap"); err != nil {
		return "", err
	}

	return <-r.addCallback(packetId, true), nil
}

func (r *rconImpl) Start() {
	go func() {
		for {
			fmt.Printf("Trying to read packet\n")
			packet := packet{}
			_, err := packet.ReadFrom(r.conn)

			switch {
			case errors.Is(err, net.ErrClosed):
				return
			case errors.Is(err, io.ErrUnexpectedEOF):
				// Can happen when connection is closed by server due to inactivity
				return
			case err != nil:
				fmt.Println("Error reading from connection:", err)
				continue
			}

			fmt.Printf(
				"Packet received; Id: %d, Type: %d, Body size: %d, Body: %s\n",
				packet.Id,
				packet.Type,
				packet.GetBodySize(),
				packet.GetBody(),
			)

			if !r.authenticated {
				switch packet.Type {
				case serverDataResponseValue:
					if packet.GetBodySize() > 0 {
						fmt.Printf(
							"Discarding non-empty data packet while not authenticated %v\n",
							packet,
						)
					}
					continue
				}
			}

			// Completion packets, identified by odd IDs, signal completeness for responses with
			// ID - 1. E.g.:
			// When receiving a packet with an even ID, we append the data to the callbacks[ID].Data.
			// When receiving a packet with an odd ID, we know that callbacks[ID - 1] is complete.
			isCompletionPacket := packet.Id%2 == 1
			callbackId := packet.Id

			if isCompletionPacket {
				callbackId--
			}

			r.callbackLock.Lock()
			callback, exists := r.callbacks[callbackId]
			if !exists {
				fmt.Printf("Callback for ID %d not registered\n", packet.Id)
				r.callbackLock.Unlock()
				continue
			}

			isComplete := false
			if callback.MaybeMultiPacket {
				if isCompletionPacket {
					isComplete = true
				} else {
					callback.Data = append(callback.Data, packet.Body...)
				}
			} else {
				callback.Data = packet.Body
				isComplete = true
			}

			if isComplete {
				callback.Channel <- string(callback.Data)
				close(callback.Channel)
				delete(r.callbacks, packet.Id)
			}
			r.callbackLock.Unlock()
		}
	}()
}

// authenticate sends a serverDataAuth request and authenticates the following requests.
func (r *rconImpl) authenticate(password string) error {
	r.authenticated = false
	packetId := r.getNextId()
	if err := r.write(serverDataAuth, packetId, password); err != nil {
		return err
	}

	select {
	case <-r.addCallback(packetId, false):
		// Login success
		r.authenticated = true
		return nil
	case <-r.addCallback(-1, false):
		// Login failure
		return fmt.Errorf("auth failed, TK change error")
	}
}

func (r *rconImpl) addCallback(id int32, maybeMultiPacket bool) chan string {
	channel := make(chan string)
	r.callbackLock.Lock()
	r.callbacks[id] = &callback{
		Channel:          channel,
		Data:             make([]byte, 0),
		MaybeMultiPacket: maybeMultiPacket,
	}
	r.callbackLock.Unlock()
	return channel
}

func (r *rconImpl) write(packetType int32, packetId int32, command string) error {
	if r.writeTimeout != 0 {
		if err := r.conn.SetWriteDeadline(time.Now().Add(r.writeTimeout)); err != nil {
			return fmt.Errorf("failed to set write deadline: %w", err)
		}
	}

	packet := newPacket(packetType, packetId, command)
	_, err := packet.WriteTo(r.conn)

	return err
}

// getNextId returns the next packet ID.
// Will always return even numbers.
// The returned number will be between startId and startId + wrapIdsAfter.
func (r *rconImpl) getNextId() int32 {
	r.idCounterLock.Lock()
	defer r.idCounterLock.Unlock()
	if r.execIdCounter < r.startId || r.execIdCounter > r.startId+wrapIdsAfter+1 {
		r.execIdCounter = r.startId
	}

	if r.execIdCounter%2 == 1 {
		r.execIdCounter++
	}

	result := int32(r.execIdCounter)
	r.execIdCounter += 2
	return result
}
