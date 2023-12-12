package squadrcon

import (
	"errors"
	"github.com/gorcon/rcon"
)

var (
	ErrNotConnected = errors.New("no connection, use .Connect to create a connection")
)

type SquadRcon struct {
	connection *rcon.Conn
}

func (r *SquadRcon) Connect(address string, password string) error {
	connection, err := rcon.Dial(address, password)
	if err != nil {
		return err
	}

	r.connection = connection

	return nil
}
