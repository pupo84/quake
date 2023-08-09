package entity

const (
	INIT_GAME = "InitGame"
	WORD      = "<world>"
)

type Games []*Game

type Game struct {
	ID             int
	Kills          int
	Players        []string
	KillsByPlayers map[string]int
	KillsByCause   map[string]int
}

func NewGame(id int) *Game {
	return &Game{
		ID:             id,
		Kills:          0,
		Players:        make([]string, 0),
		KillsByPlayers: make(map[string]int),
		KillsByCause:   make(map[string]int),
	}
}

func (g *Game) AddPlayer(player string) {
	for _, p := range g.Players {
		if p == player {
			return
		}
	}
	g.Players = append(g.Players, player)
}

func (g *Game) AddPlayerStats(player, cause string, value int) {
	g.KillsByPlayers[player] += value
	g.KillsByCause[cause]++
	g.Kills++
}

func (g *Game) HasStats() bool {
	return len(g.Players) > 0
}
