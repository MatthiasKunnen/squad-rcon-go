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
	Data string
}

type rconImpl struct {
	authenticated bool
	callbackLock  sync.Mutex
	callbacks     map[int32]*callback
	conn          net.Conn
	dialTimeout   time.Duration
	execIdCounter int
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

	return <-r.addCallback(packetId), nil
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

			r.callbackLock.Lock()
			callback, exists := r.callbacks[packet.Id]
			if !exists {
				fmt.Printf("Callback for ID %d not registered\n", packet.Id)
				r.callbackLock.Unlock()
				continue
			}

			callback.Channel <- packet.GetBody()
			close(callback.Channel)
			delete(r.callbacks, packet.Id)
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
	case <-r.addCallback(packetId):
		// Login success
		r.authenticated = true
		return nil
	case <-r.addCallback(-1):
		// Login failure
		return fmt.Errorf("auth failed, TK change error")
	}
}

func (r *rconImpl) addCallback(id int32) chan string {
	channel := make(chan string)
	r.callbackLock.Lock()
	r.callbacks[id] = &callback{
		Channel: channel,
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

func (r *rconImpl) getNextId() int32 {
	r.execIdCounter++
	if r.execIdCounter > r.startId+wrapIdsAfter {
		r.execIdCounter = r.startId
	}

	return int32(r.execIdCounter)
}
