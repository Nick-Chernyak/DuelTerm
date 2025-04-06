package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"duelterm/pkg/common"
)

const (
	ArenaWidth   = 20
	ArenaHeight  = 10
	TickRate     = 100 * time.Millisecond
	MaxHP        = 3
	ShotCooldown = 20
)

type Client struct {
	Conn      net.Conn
	Enc       *json.Encoder
	Dec       *json.Decoder
	Player    common.PlayerState
	Direction string
	Cooldown  int
}

var (
	clients     [2]*Client
	projectiles []common.Projectile
	lock        sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Server] Waiting for 2 players...")

	for i := 0; i < 2; i++ {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		clients[i] = initClient(conn, i)
		go handleClientInput(i)
		fmt.Printf("[Server] Player %d connected\n", i+1)
	}

	gameLoop()
}

func initClient(conn net.Conn, index int) *Client {
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	name := fmt.Sprintf("Player%d", index+1)
	ch := '@'
	if index == 1 {
		ch = '&'
	}
	return &Client{
		Conn:      conn,
		Enc:       enc,
		Dec:       dec,
		Direction: "right",
		Player: common.PlayerState{
			Name: name,
			HP:   MaxHP,
			X:    index * (ArenaWidth - 1),
			Y:    ArenaHeight / 2,
			Char: ch,
		},
	}
}

func handleClientInput(index int) {
	client := clients[index]
	for {
		var msg common.ActionMessage
		if err := client.Dec.Decode(&msg); err != nil {
			continue
		}
		lock.Lock()
		handleAction(client, msg)
		lock.Unlock()
	}
}

func handleAction(c *Client, msg common.ActionMessage) {
	switch msg.Action {
	case "move":
		dir := msg.Direction
		d := common.DirectionVectors[dir]
		newX := c.Player.X + d[0]
		newY := c.Player.Y + d[1]
		if newX >= 0 && newX < ArenaWidth {
			c.Player.X = newX
		}
		if newY >= 0 && newY < ArenaHeight {
			c.Player.Y = newY
		}
		c.Direction = dir
	case "attack":
		if c.Cooldown <= 0 {
			vec := common.DirectionVectors[c.Direction]
			p := common.Projectile{
				X:     c.Player.X + vec[0],
				Y:     c.Player.Y + vec[1],
				DirX:  vec[0],
				DirY:  vec[1],
				Owner: c.Player.Name,
			}
			projectiles = append(projectiles, p)
			c.Cooldown = ShotCooldown
		}
	}
}

func gameLoop() {
	ticker := time.NewTicker(TickRate)
	for range ticker.C {
		lock.Lock()
		updateProjectiles()
		for _, c := range clients {
			if c.Cooldown > 0 {
				c.Cooldown--
			}
		}
		gs := common.GameState{
			ArenaWidth:  ArenaWidth,
			ArenaHeight: ArenaHeight,
			Players:     []common.PlayerState{clients[0].Player, clients[1].Player},
			Projectiles: projectiles,
		}
		if clients[0].Player.HP <= 0 {
			gs.Message = fmt.Sprintf("%s wins!", clients[1].Player.Name)
		} else if clients[1].Player.HP <= 0 {
			gs.Message = fmt.Sprintf("%s wins!", clients[0].Player.Name)
		}
		for _, c := range clients {
			_ = c.Enc.Encode(gs)
		}
		lock.Unlock()
	}
}

func updateProjectiles() {
	newProjectiles := make([]common.Projectile, 0)
	for _, p := range projectiles {
		p.X += p.DirX
		p.Y += p.DirY

		if p.X < 0 || p.X >= ArenaWidth || p.Y < 0 || p.Y >= ArenaHeight {
			continue
		}

		for i := range clients {
			target := &clients[i].Player
			if target.Name != p.Owner && target.HP > 0 && target.X == p.X && target.Y == p.Y {
				target.HP--
				fmt.Printf("[HIT] %s hit %s! %d HP left\n", p.Owner, target.Name, target.HP)
				goto skip
			}
		}

		newProjectiles = append(newProjectiles, p)
	skip:
	}
	projectiles = newProjectiles
}
