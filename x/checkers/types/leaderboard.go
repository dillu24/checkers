package types

import "fmt"

func (leaderboard Leaderboard) Validate() error {
	// Check for duplicated player addresses in winners
	winnerInfoIndexMap := make(map[string]struct{})

	for _, elem := range leaderboard.Winners {
		index := string(PlayerInfoKey(elem.PlayerAddress))
		if _, ok := winnerInfoIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for winner")
		}
		winnerInfoIndexMap[index] = struct{}{}
	}
	return nil
}
