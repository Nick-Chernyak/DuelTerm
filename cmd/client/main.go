package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"duelterm/pkg/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	conn      net.Conn
	enc       *json.Encoder
	dec       *json.Decoder
	gameState common.GameState
	connected bool
	error     error
}

func (m model) Init() tea.Cmd {
	return readGameState(m.dec)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.connected {
			return m, nil
		}
		switch msg.String() {
		case "w":
			m.enc.Encode(common.ActionMessage{Action: "move", Direction: "up"})
		case "s":
			m.enc.Encode(common.ActionMessage{Action: "move", Direction: "down"})
		case "a":
			m.enc.Encode(common.ActionMessage{Action: "move", Direction: "left"})
		case "d":
			m.enc.Encode(common.ActionMessage{Action: "move", Direction: "right"})
		case "m":
			m.enc.Encode(common.ActionMessage{Action: "attack"})
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case common.GameState:
		m.gameState = msg
		m.connected = true
		return m, readGameState(m.dec)
	case error:
		m.error = msg
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if m.error != nil {
		return "Error: " + m.error.Error()
	}
	if !m.connected {
		return "Подключение к серверу..."
	}
	return renderGame(m.gameState)
}

func main() {
	addr := "localhost:8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	m := model{
		conn: conn,
		enc:  enc,
		dec:  dec,
	}
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func readGameState(dec *json.Decoder) tea.Cmd {
	return func() tea.Msg {
		var state common.GameState
		if err := dec.Decode(&state); err != nil {
			return err
		}
		return state
	}
}

func renderGame(state common.GameState) string {
	var out string
	arena := make([][]rune, state.ArenaHeight)
	for y := range arena {
		arena[y] = make([]rune, state.ArenaWidth)
		for x := range arena[y] {
			arena[y][x] = '.'
		}
	}
	for _, p := range state.Players {
		if p.HP > 0 {
			arena[p.Y][p.X] = p.Char
		}
	}

	for _, proj := range state.Projectiles {
		if proj.Y >= 0 && proj.Y < len(arena) && proj.X >= 0 && proj.X < len(arena[0]) {
			arena[proj.Y][proj.X] = getBulletChar(bulletIndex)
			bulletIndex++
			if bulletIndex >= len(bulletChars) {
				bulletIndex = 0
			}
		}
	}

	for _, row := range arena {
		for _, c := range row {
			out += string(c)
		}
		out += "\n"
	}
	for _, p := range state.Players {
		bar := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%s: %d HP", p.Name, p.HP))
		out += bar + "\n"
	}
	if state.Message != "" {
		out += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true).Render(state.Message) + "\n"
	}
	return out
}

var bulletChars = []rune{'з', 'а', 'л'}

var bulletIndex = 0

func getBulletChar(index int) rune {
	return bulletChars[index%len(bulletChars)]
}
