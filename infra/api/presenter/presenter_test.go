package presenter_test

import (
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/pupo84/quake/domain/entity"
	"github.com/pupo84/quake/infra/api/presenter"
	"github.com/stretchr/testify/assert"
)

func TestNewGameResponse(t *testing.T) {
	games := []*entity.Game{
		{
			ID:             1,
			Kills:          10,
			Players:        []string{"player1", "player2"},
			KillsByPlayers: map[string]int{"player1": 5, "player2": 5},
			KillsByCause:   map[string]int{"MOD_SHOTGUN": 10},
		},
		{
			ID:             2,
			Kills:          5,
			Players:        []string{"player1", "player3"},
			KillsByPlayers: map[string]int{"player1": 2, "player3": 3},
			KillsByCause:   map[string]int{"MOD_ROCKET": 5},
		},
	}

	expectedResponse := orderedmap.New()
	expectedResponse.Set("game_1", &presenter.GameBody{
		TotalKills:   10,
		Players:      []string{"player1", "player2"},
		Kills:        map[string]int{"player1": 5, "player2": 5},
		KillsByMeans: map[string]int{"MOD_SHOTGUN": 10},
	})
	expectedResponse.Set("game_2", &presenter.GameBody{
		TotalKills:   5,
		Players:      []string{"player1", "player3"},
		Kills:        map[string]int{"player1": 2, "player3": 3},
		KillsByMeans: map[string]int{"MOD_ROCKET": 5},
	})

	response := presenter.NewGameResponse(games)

	for _, key := range expectedResponse.Keys() {
		expectedValue, _ := expectedResponse.Get(key)
		value, ok := response.Get(key)

		if !ok {
			t.Errorf("Expected key %s not found in response", key)
		}

		assert.Equal(t, expectedValue, value)
	}
}
