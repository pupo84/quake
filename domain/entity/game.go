package entity

// Constants for game initialization and world string
const (
	INIT_GAME = "InitGame"
	WORD      = "<world>"
)

// Games is a slice of Game pointers
type Games []*Game

// Game represents a Quake game
type Game struct {
	ID             int            // ID of the game
	Kills          int            // Total number of kills in the game
	Players        []string       // List of players in the game
	KillsByPlayers map[string]int // Map of kills by each player
	KillsByCause   map[string]int // Map of kills by each cause
}

// NewGame creates a new Game instance with the given ID
func NewGame(id int) *Game {
	return &Game{
		ID:             id,
		Kills:          0,
		Players:        make([]string, 0),
		KillsByPlayers: make(map[string]int),
		KillsByCause:   make(map[string]int),
	}
}

// AddPlayer adds a player to the game's list of players
func (g *Game) AddPlayer(player string) {
	for _, p := range g.Players {
		if p == player {
			return
		}
	}
	g.Players = append(g.Players, player)
}

// AddPlayerStats adds stats for a player in the game
func (g *Game) AddPlayerStats(player, cause string, value int) {
	g.KillsByPlayers[player] += value
	g.KillsByCause[cause]++
	g.Kills++
}

// HasStats returns true if the game has any player stats
func (g *Game) HasStats() bool {
	return len(g.Players) > 0
}
