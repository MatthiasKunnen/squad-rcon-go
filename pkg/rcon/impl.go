package rcon

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
	ErrCommandEmpty      = errors.New("command is empty")
	ErrIncorrectPassword = errors.New("RCON password is incorrect")
	ErrNotAuthenticated  = errors.New("not authenticated")
)

type rconResponse struct {
	Body []byte
}

type callback struct {
	// The channel which will be passed the data as responded by RCON.
	Channel chan string

	// The aggregated data.
	Data []byte
}

type rconImpl struct {
	// Whether conn is authenticated.
	authenticated bool

	// Lock to be used before accessing the callbacks map.
	callbackLock sync.Mutex

	// Map of callbacks, key: packet ID, value callback to call when response is received.
	callbacks map[int32]*callback

	// The command to send after every Execute.
	// Responses to this confirmation command confirm that all responses to the Execute command have
	// been received.
	confirmationCommand string

	// TCP connection with RCON
	conn net.Conn

	// The time after which connecting is aborted.
	dialTimeout time.Duration

	// The next packet ID to use
	execIdCounter int

	// Lock needed before accessing execIdCounter.
	idCounterLock sync.Mutex

	// The first packet ID that can be used.
	startId int

	// The time to wait for writing to complete before aborting.
	writeTimeout time.Duration
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
	if err := r.write(serverDataExecCommand, packetId+1, r.confirmationCommand); err != nil {
		return "", err
	}

	return <-r.addCallback(packetId), nil
}

func (r *rconImpl) start() {
	go func() {
		for {
			if err := r.handleIncomingPacket(); err != nil {
				// todo, set up reconnect on error
				fmt.Printf("Encountered an error processing packets. Stopping. %v\n", err)
				return
			}
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

	emptyPacket := packet{}
	_, err := emptyPacket.ReadFrom(r.conn)

	// Squad's RCON implementation closes the connection on failed authentication
	switch {
	case errors.Is(err, net.ErrClosed):
		return err
	case errors.Is(err, io.EOF):
		return ErrIncorrectPassword
	case errors.Is(err, io.ErrUnexpectedEOF):
		// Can happen when connection is closed by server due to inactivity
		return err
	case err != nil:
		fmt.Println("Error reading from connection:", err)
		return err
	}

	authResultPacket := packet{}
	_, err = authResultPacket.ReadFrom(r.conn)

	if err != nil {
		return err
	}

	if authResultPacket.Id != packetId {
		return fmt.Errorf("unexpected ID in auth response. Got %d, expected %d", authResultPacket.Id, packetId)
	}

	r.authenticated = true
	return nil
}

func (r *rconImpl) addCallback(id int32) chan string {
	channel := make(chan string)
	r.callbackLock.Lock()
	r.callbacks[id] = &callback{
		Channel: channel,
		Data:    make([]byte, 0),
	}
	r.callbackLock.Unlock()
	return channel
}

// handleIncomingPacket processes all incoming packets after authentication.
// It assumes that responses are multi-packet and are followed by a _confirmation command_.
func (r *rconImpl) handleIncomingPacket() error {
	if !r.authenticated {
		return ErrNotAuthenticated
	}

	fmt.Printf("Trying to read packet\n")
	packet := packet{}
	_, err := packet.ReadFrom(r.conn)

	switch {
	case errors.Is(err, net.ErrClosed):
		return err
	case errors.Is(err, io.ErrUnexpectedEOF):
		// Can happen when connection is closed by server due to inactivity
		return err
	case err != nil:
		return err
	}

	fmt.Printf(
		"Packet received; Id: %d, Type: %d, Body size: %d, Body: %s\n",
		packet.Id,
		packet.Type,
		packet.GetBodySize(),
		packet.GetBody(),
	)

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
		return nil
	}

	if isCompletionPacket {
		callback.Channel <- string(callback.Data)
		close(callback.Channel)
		delete(r.callbacks, packet.Id)
	} else {
		callback.Data = append(callback.Data, packet.Body...)
	}

	r.callbackLock.Unlock()
	return nil
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
