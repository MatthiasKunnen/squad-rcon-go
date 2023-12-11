# RCON `ListPlayers`

## Notes
- Users first join without their name being prefixed
- The command seems to be able to return a string that is not a playerlist, e.g.
   ```
   Jon (Steam ID: 76561197999957991) has created Squad 1 (Squad Name: TESTICLES) on United States Army
   ```

```
2023-11-30-19:41:34
----- Active Players -----
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:35
----- Active Players -----
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:36
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:37
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:38
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:39
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:40
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:41
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:42
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:43
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:44
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:45
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:46
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:47
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:48
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:49
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:50
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: N/A | Is Leader: False | Role: USA_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:51
Jon (Steam ID: 76561197999957991) has created Squad 1 (Squad Name: TESTICLES) on United States Army
2023-11-30-19:41:52
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: 1 | Is Leader: True | Role: USA_SL_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-19:41:53
```

### ListPlayers respond is not a list
```
2023-11-30-23:56:27
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: 1 | Is Leader: True | Role: USA_SL_01
ID: 1 | SteamID: 76561197989362395 | Name: ✯RAIDR✯creaman | Team ID: 2 | Squad ID: N/A | Is Leader: False | Role: INS_Rifleman_01
----- Recently Disconnected Players [Max of 15] -----
2023-11-30-23:56:28
creaman (Steam ID: 76561197989362395) has created Squad 1 (Squad Name: ILU WAFFLES) on Insurgent Forces
2023-11-30-23:56:29
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: 1 | Is Leader: True | Role: USA_SL_01
ID: 1 | SteamID: 76561197989362395 | Name: ✯RAIDR✯creaman | Team ID: 2 | Squad ID: 1 | Is Leader: True | Role: INS_SL_01
----- Recently Disconnected Players [Max of 15] -----
```

### Disconnect
```
2023-11-30-23:59:23
----- Active Players -----
ID: 0 | SteamID: 76561197999957991 | Name: ✯RAIDR✯Jon | Team ID: 1 | Squad ID: 1 | Is Leader: True | Role: USA_Pilot_01
----- Recently Disconnected Players [Max of 15] -----
ID: 1 | SteamID: 76561197989362395 | Since Disconnect: 00m.04s | Name: creaman
```
