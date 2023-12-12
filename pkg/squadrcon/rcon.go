package squadrcon

import (
	"fmt"
	"net"
	"time"
)

type Rcon interface {
	Close() error
	Execute(command string) (string, error)
	Start()
}

type RconSettings struct {
	DialTimeout time.Duration

	// PacketIdStart contains the first packet ID that will be used. Change it when multiple rcon
	// connections are used. E.g. SquadJS uses ID 1 and 2, so these IDs shouldn't be used to prevent
	// conflicts.
	PacketIdStart int32

	WriteTimeout time.Duration
}

// Connect connects to the RCON server and authenticates.
func Connect(address string, password string, settings RconSettings) (Rcon, error) {
	client := &rconImpl{
		dialTimeout:   5 * time.Second,
		writeTimeout:  5 * time.Second,
		callbacks:     make(map[int32]callback),
		execIdCounter: 10000,
		startId:       10000,
	}

	if settings.DialTimeout > 0 {
		client.dialTimeout = settings.DialTimeout
	}

	if settings.WriteTimeout > 0 {
		client.writeTimeout = settings.WriteTimeout
	}

	conn, err := net.DialTimeout("tcp", address, settings.DialTimeout)
	if err != nil {
		// Failed to open TCP connection to the server.
		return nil, fmt.Errorf("failed to connect to rcon server on %s: %w", address, err)
	}
	client.conn = conn
	client.Start()

	if err := client.authenticate(password); err != nil {
		if err2 := client.Close(); err2 != nil {
			return client, fmt.Errorf(
				"failed to close connection: %v. Previous error: %v",
				err2,
				err,
			)
		}

		return client, fmt.Errorf("failed to authenticate rcon connection: %w", err)
	}

	return client, nil
}
