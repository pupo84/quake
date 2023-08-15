package presenter

import (
	"fmt"

	"github.com/iancoleman/orderedmap"
	"github.com/pupo84/quake/domain/entity"
)

// Response is the response body for a game
type Response map[string]*GameBody

// GameBody is the response body for a game
type GameBody struct {
	TotalKills   int            `json:"total_kills"`
	Players      []string       `json:"players"`
	Kills        map[string]int `json:"kills,omitempty"`
	KillsByMeans map[string]int `json:"kills_by_means,omitempty"`
}

// NewGameResponse returns a new game response
func NewGameResponse(games []*entity.Game) *orderedmap.OrderedMap {
	response := orderedmap.New()

	for _, game := range games {
		body := &GameBody{
			TotalKills:   game.Kills,
			Players:      game.Players,
			Kills:        game.KillsByPlayers,
			KillsByMeans: game.KillsByCause,
		}

		gameID := fmt.Sprintf("game_%d", game.ID)
		response.Set(gameID, body)
	}

	return response
}
