package v1tov2

import "github.com/alice/checkers/x/checkers/types"

const (
	StoredGameChunkSize = 1_000
	PlayerInfoChunkSize = types.LeaderboardWinnerLength * 2
)
