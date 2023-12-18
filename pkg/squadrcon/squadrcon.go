package squadrcon

import (
	"errors"
	"squad-rcon-go/pkg/rcon"
	"time"
)

var (
	ErrNotConnected = errors.New("no connection, use .Connect to create a connection")
)

type SquadRcon struct {
	rcon rcon.Rcon
}

type Settings struct {
	DialTimeout time.Duration

	// PacketIdStart contains the first packet ID that will be used. Change it when multiple rcon
	// connections are used. E.g. SquadJS uses ID 1 and 2, so these IDs shouldn't be used to prevent
	// conflicts.
	PacketIdStart int32

	WriteTimeout time.Duration
}

func Connect(address string, password string, settings Settings) (rcon.Rcon, error) {
	rc, err := rcon.Connect(address, password, rcon.Settings{
		DialTimeout:   settings.DialTimeout,
		PacketIdStart: settings.PacketIdStart,
		WriteTimeout:  settings.WriteTimeout,
	})
	if err != nil {
		return nil, err
	}

	squadRcon := &SquadRcon{
		rcon: rc,
	}

	return squadRcon, nil
}

func (r *SquadRcon) Close() error {
	return r.rcon.Close()
}

func (r *SquadRcon) Execute(command string) (string, error) {
	return r.rcon.Execute(command)
}
