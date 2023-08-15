package entity

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame(1)

	if game.ID != 1 {
		t.Errorf("Expected game ID to be 1, but got %d", game.ID)
	}

	if game.Kills != 0 {
		t.Errorf("Expected game Kills to be 0, but got %d", game.Kills)
	}

	if len(game.Players) != 0 {
		t.Errorf("Expected game Players to be empty, but got %v", game.Players)
	}

	if len(game.KillsByPlayers) != 0 {
		t.Errorf("Expected game KillsByPlayers to be empty, but got %v", game.KillsByPlayers)
	}

	if len(game.KillsByCause) != 0 {
		t.Errorf("Expected game KillsByCause to be empty, but got %v", game.KillsByCause)
	}
}

func TestAddPlayer(t *testing.T) {
	game := NewGame(1)

	game.AddPlayer("player1")
	game.AddPlayer("player2")
	game.AddPlayer("player1")

	if len(game.Players) != 2 {
		t.Errorf("Expected game Players to have 2 elements, but got %v", game.Players)
	}
}

func TestAddPlayerStats(t *testing.T) {
	game := NewGame(1)

	game.AddPlayerStats("player1", "cause1", 1)
	game.AddPlayerStats("player2", "cause2", 1)
	game.AddPlayerStats("player2", "cause1", 1)

	if game.Kills == 2 {
		t.Errorf("Expected game Kills to be 4, but got %d", game.Kills)
	}

	if game.KillsByPlayers["player1"] != 1 {
		t.Errorf("Expected game KillsByPlayers[player1] to be 1, but got %d", game.KillsByPlayers["player1"])
	}

	if game.KillsByPlayers["player2"] != 2 {
		t.Errorf("Expected game KillsByPlayers[player2] to be 2, but got %d", game.KillsByPlayers["player2"])
	}

	if game.KillsByCause["cause1"] != 2 {
		t.Errorf("Expected game KillsByCause[cause1] to be 2, but got %d", game.KillsByCause["cause1"])
	}

	if game.KillsByCause["cause2"] != 1 {
		t.Errorf("Expected game KillsByCause[cause2] to be 1, but got %d", game.KillsByCause["cause2"])
	}
}

func TestHasStats(t *testing.T) {
	game := NewGame(1)

	if game.HasStats() {
		t.Errorf("Expected game HasStats to be false, but got true")
	}

	game.AddPlayer("player1")

	if !game.HasStats() {
		t.Errorf("Expected game HasStats to be true, but got false")
	}
}
