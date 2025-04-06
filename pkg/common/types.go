package common

var DirectionVectors = map[string][2]int{
	"up":    {0, -1},
	"down":  {0, 1},
	"left":  {-1, 0},
	"right": {1, 0},
}

type ActionMessage struct {
	Action    string `json:"action"`
	Direction string `json:"direction"`
}

type PlayerState struct {
	Name string `json:"name"`
	HP   int    `json:"hp"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Char rune   `json:"char"`
}

type Projectile struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	DirX  int    `json:"dir_x"`
	DirY  int    `json:"dir_y"`
	Owner string `json:"owner"`
}

type GameState struct {
	ArenaWidth  int           `json:"arena_width"`
	ArenaHeight int           `json:"arena_height"`
	Players     []PlayerState `json:"players"`
	Message     string        `json:"message"`
	Projectiles []Projectile  `json:"projectiles"`
}
