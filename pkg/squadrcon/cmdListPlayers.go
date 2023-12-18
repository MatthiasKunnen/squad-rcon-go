package squadrcon

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"squad-rcon-go/pkg/rcon"
	"strconv"
	"strings"
	"time"
)

type ActivePlayer struct {
	MatchId int

	SteamId string

	Name string

	// TeamId contains the team index
	TeamIndex int

	// SquadId contains the 1-indexed index of the player's squad. If zero, then the player is not
	// part of a squad.
	SquadIndex int

	IsSquadLead bool

	Kit string
}

type DisconnectedPlayer struct {
	MatchId int

	SteamId string

	Name string

	DisconnectTime time.Time
}

type PlayerList struct {
	ActivePlayers       []ActivePlayer
	DisconnectedPlayers []DisconnectedPlayer
}

var (
	ErrResponseIsNotPlayerList = errors.New("response returned from rcon is not a player list")
)

func ListPlayers(rcon rcon.Rcon) (PlayerList, error) {
	response, err := rcon.Execute("ListPlayers")
	if err != nil {
		return PlayerList{}, err
	}

	list, err := ParsePlayersList(response)

	//if errors.Is(err, ErrResponseIsNotPlayerList) {
	//	time.Sleep(500 * time.Millisecond)
	//	return ListPlayers(rcon)
	//}

	if err != nil {
		return PlayerList{}, err
	}

	return list, nil
}

type parsePlayerListState int

const (
	LookingForActivePlayerHeader parsePlayerListState = iota
	ReadingActivePlayers
	ReadingDisconnectedPlayers
)

const activePlayersHeader = "----- Active Players -----"
const disconnectedPlayersPrefix = "----- Recently Disconnected Players "

var playerListActivePlayerRegex = regexp.MustCompile(`^ID: (\d+) \| SteamID: (\d+) \| Name: (.+) \| Team ID: (\d+) \| Squad ID: ([^|]+) \| Is Leader: (\w+) \| Role: ([^|]+)$`)

const (
	_ = iota
	activePlayerMatchId
	activePlayerSteamId
	activePlayerName
	activePlayerTeamIndex
	activePlayerSquadIndex
	activePlayerIsSquadLead
	activePlayerRole
)

var playerListDisconnectedPlayerRegex = regexp.MustCompile(`^ID: (\d+) \| SteamID: (\d+) \| Since Disconnect: (\d+)m.(\d+)s \| Name: (.+)$`)

const (
	_ = iota
	disconnectedPlayerMatchIdIndex
	disconnectedPlayerSteamIdIndex
	disconnectedPlayerMinutesIndex
	disconnectedPlayerSecondsIndex
	disconnectedPlayerNameIndex
)

func ParsePlayersList(playerListString string) (PlayerList, error) {
	if playerListString == "" {
		return PlayerList{}, nil
	}

	if !strings.HasPrefix(playerListString, activePlayersHeader) {
		return PlayerList{}, ErrResponseIsNotPlayerList
	}

	var playerList PlayerList
	var state parsePlayerListState
	var errs []error
	scanner := bufio.NewScanner(strings.NewReader(playerListString))

	for scanner.Scan() {
		var line = scanner.Text()

		switch state {
		case LookingForActivePlayerHeader:
			if !strings.HasPrefix(line, activePlayersHeader) {
				errs = append(
					errs,
					fmt.Errorf("expected \"%s\", got \"%s\"", activePlayersHeader, line),
				)
				continue
			}
			state = ReadingActivePlayers
		case ReadingActivePlayers:
			if strings.HasPrefix(line, disconnectedPlayersPrefix) {
				state = ReadingDisconnectedPlayers
				continue
			}

			matches := playerListActivePlayerRegex.FindStringSubmatch(line)
			fmt.Printf("%q\n", matches)

			if matches == nil {
				errs = append(
					errs,
					fmt.Errorf("line cannot be parsed as an active player: %s", line),
				)
				continue
			}

			var playerMatchIdString = matches[activePlayerMatchId]
			playerMatchId, err := strconv.Atoi(playerMatchIdString)
			if err != nil {
				errs = append(
					errs,
					fmt.Errorf("could not parse \"%s\" as player match ID", playerMatchIdString),
				)
				continue
			}

			var playerSquadIndexString = matches[activePlayerSquadIndex]
			var playerSquadIndex = 0
			if playerSquadIndexString != "N/A" {
				playerSquadIndex, err = strconv.Atoi(playerSquadIndexString)
				if err != nil {
					errs = append(
						errs,
						fmt.Errorf(
							"could not parse player squad index \"%s\"",
							playerSquadIndexString,
						),
					)
					continue
				}
			}

			var playerTeamIndexString = matches[activePlayerTeamIndex]
			playerTeamIndex, err := strconv.Atoi(playerTeamIndexString)
			if err != nil {
				errs = append(
					errs,
					fmt.Errorf(
						"could not parse player team index \"%s\"",
						playerTeamIndexString,
					),
				)
				continue
			}

			playerList.ActivePlayers = append(playerList.ActivePlayers, ActivePlayer{
				IsSquadLead: matches[activePlayerIsSquadLead] == "True",
				Kit:         matches[activePlayerRole],
				MatchId:     playerMatchId,
				Name:        matches[activePlayerName],
				SquadIndex:  playerSquadIndex,
				SteamId:     matches[activePlayerSteamId],
				TeamIndex:   playerTeamIndex,
			})
		case ReadingDisconnectedPlayers:
			matches := playerListDisconnectedPlayerRegex.FindStringSubmatch(line)

			if matches == nil {
				errs = append(
					errs,
					fmt.Errorf("line cannot be parsed as a disconnected player: %s", line),
				)
				continue
			}

			var disconnectedPlayerMatchIdString = matches[disconnectedPlayerMatchIdIndex]
			playerMatchId, err := strconv.Atoi(disconnectedPlayerMatchIdString)
			if err != nil {
				errs = append(
					errs,
					fmt.Errorf(
						"could not parse disconnected player match id \"%s\"",
						disconnectedPlayerMatchIdString,
					),
				)
				continue
			}

			var disconnectedPlayerMinutesString = matches[disconnectedPlayerMinutesIndex]
			disconnectedMinutes, err := strconv.Atoi(disconnectedPlayerMinutesString)
			if err != nil {
				errs = append(
					errs,
					fmt.Errorf(
						"could not parse disconnected player minutes \"%s\"",
						disconnectedPlayerMatchIdString,
					),
				)
				continue
			}

			var disconnectedPlayerSecondsString = matches[disconnectedPlayerSecondsIndex]
			disconnectedSeconds, err := strconv.Atoi(disconnectedPlayerSecondsString)
			if err != nil {
				errs = append(
					errs,
					fmt.Errorf(
						"could not parse disconnected player seconds \"%s\"",
						disconnectedPlayerSecondsString,
					),
				)
				continue
			}
			var disconnectTime = time.Now().Add(-time.Duration(disconnectedSeconds+(disconnectedMinutes*60)) * time.Second)

			playerList.DisconnectedPlayers = append(playerList.DisconnectedPlayers, DisconnectedPlayer{
				MatchId:        playerMatchId,
				SteamId:        matches[disconnectedPlayerSteamIdIndex],
				Name:           matches[disconnectedPlayerNameIndex],
				DisconnectTime: disconnectTime,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return playerList, errors.Join(errs...)
	}

	return playerList, nil
}
