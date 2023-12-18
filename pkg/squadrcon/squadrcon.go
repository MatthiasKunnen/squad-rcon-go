package squadrcon

import (
	"errors"
	"squad-rcon-go/pkg/rcon"
)

var (
	ErrNotConnected = errors.New("no connection, use .Connect to create a connection")
)

type SquadRcon struct {
	rcon rcon.Rcon
}

func Connect(address string, password string, settings rcon.Settings) (rcon.Rcon, error) {
	rc, err := rcon.Connect(address, password, settings)
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
