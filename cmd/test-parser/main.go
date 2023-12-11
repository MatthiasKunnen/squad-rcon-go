package main

import (
	"fmt"
	"squad-rcon-go/pkg/squadrcon"
)

var examples = []string{
	`----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: 1 | Is Leader: True | Role: USA_SL_01
ID: 1 | SteamID: 76561197989362395 | Name: ✯RAIDR✯creaman | Team ID: 2 | Squad ID: N/A | Is Leader: False | Role: INS_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
`,
	`----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----`,
	`----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: 1 | Is Leader: True | Role: USA_Pilot_01
----- Recently Disconnected Players [Max of 15] -----
ID: 1 | SteamID: 76561197989362395 | Since Disconnect: 00m.04s | Name: creaman`,
	`creaman (Steam ID: 76561197989362395) has created Squad 1 (Squad Name: ILU WAFFLES) on Insurgent Forces`,
	`----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: J👨‍❤️‍👨🍛♟ | Team ID: 1 | Squad ID: 5 | Is Leader: False | Role: USA_Rifleman_01
ID: 0 | SteamID: 76561197999957991 | Name: Juehn | Team ID: 1 | Squad ID: 555555555555555555555555555555555554 | Is Leader: False | Role: USA_Rifleman_01
ID: 0 | SteamID: 76561197999957991 | Name: Juehn | Team ID: 1 | Squad ID: 154154546486864684864868868648686468 | Is Leader: False | Role: USA_Rifleman_01
ID: 0 | SteamID: 76561197999957991 | Name: Juehn | Team ID: 1 | Squad ID: 987987879879879879897897878897897897 | Is Leader: False | Role: USA_Rifleman_01
ID: 0 | SteamID: 76561197999957991 | Name: Juehn | Team ID: 1 | Team ID: 1 | Squad ID: 4 | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----`,
}

func main() {
	list, err := squadrcon.ParsePlayersList(examples[4])
	fmt.Println(err)
	fmt.Println(list)
}
