package consumer

import (
	"testing"

	"github.com/romangurevitch/go-training/internal/basics/interface/producerside/producer"
)

// mockStrategy must satisfy producer.Strategy — an interface owned by the producer.
// Notice we're forced to implement SomeOtherUnrelatedFunction even though the
// consumer never calls it. This is the cost of producer-side interfaces.
type mockStrategy struct {
	cmd producer.Command
}

func (m *mockStrategy) Play() producer.Command {
	return m.cmd
}

// SomeOtherUnrelatedFunction must be implemented to satisfy producer.Strategy,
// even though GameServer never calls it. This coupling is the key contrast
// with consumer-side interfaces.
func (m *mockStrategy) SomeOtherUnrelatedFunction() {}

func TestGameServer_StartGame(t *testing.T) {
	mock := &mockStrategy{cmd: producer.Forward}
	gs := NewGameServer(mock)
	gs.StartGame()
}
