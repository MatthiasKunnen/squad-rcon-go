package rcon

import (
	"fmt"
	"net"
	"time"
)

type Rcon interface {
	Close() error
	Execute(command string) (string, error)
}

type Settings struct {
	// The RCON command that is sent after every Execute.
	// Used to detect whether all responses to the execute command are received.
	// Best if the response is as small as possible.
	ConfirmationCommand string

	DialTimeout time.Duration

	// PacketIdStart contains the first packet ID that will be used. Change it when multiple rcon
	// connections are used. E.g. SquadJS uses ID 1 and 2, so these IDs shouldn't be used to prevent
	// conflicts.
	PacketIdStart int32

	WriteTimeout time.Duration
}

// Connect connects to the RCON server and authenticates.
func Connect(address string, password string, settings Settings) (Rcon, error) {
	client := &rconImpl{
		confirmationCommand: settings.ConfirmationCommand,
		dialTimeout:         5 * time.Second,
		writeTimeout:        5 * time.Second,
		callbacks:           make(map[int32]*callback),
		execIdCounter:       10000,
		startId:             10000,
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

	if err := client.authenticate(password); err != nil {
		if err2 := client.Close(); err2 != nil {
			return client, fmt.Errorf(
				"failed to close connection: %w. Previous error: %w",
				err2,
				err,
			)
		}

		return client, fmt.Errorf("failed to authenticate rcon connection: %w", err)
	}

	client.start()

	return client, nil
}
