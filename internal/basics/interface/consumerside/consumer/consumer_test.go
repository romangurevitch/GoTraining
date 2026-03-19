package consumer

import (
	"testing"

	"github.com/romangurevitch/go-training/internal/basics/interface/consumerside/producer"
)

// mockStrategy is a local mock — it satisfies the Strategy interface defined
// right here in this package. No producer import is needed to create a test double.
// This is the key benefit of consumer-side interfaces: easy, local mocking.
type mockStrategy struct {
	cmd producer.Command
}

func (m *mockStrategy) Play() producer.Command {
	return m.cmd
}

func TestGameServer_StartGame(t *testing.T) {
	mock := &mockStrategy{cmd: producer.Forward}
	gs := NewGameServer(mock)
	// StartGame prints the command — verify it doesn't panic with our mock.
	gs.StartGame()
}

func TestGameServer_UsesStrategyResult(t *testing.T) {
	commands := []producer.Command{producer.Forward, producer.Left, producer.Right, producer.Shoot}
	for _, cmd := range commands {
		t.Run(string(cmd), func(t *testing.T) {
			mock := &mockStrategy{cmd: cmd}
			gs := NewGameServer(mock)
			gs.StartGame() // each command is valid and non-panicking
		})
	}
}
