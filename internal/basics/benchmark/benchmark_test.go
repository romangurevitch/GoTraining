package benchmark

import (
	"math/big"
	"testing"
)

// run: go test -bench=. -benchmem
func BenchmarkTest2(b *testing.B)   { benchmarkTest(2, b) }
func BenchmarkTest4(b *testing.B)   { benchmarkTest(4, b) }
func BenchmarkTest8(b *testing.B)   { benchmarkTest(8, b) }
func BenchmarkTest16(b *testing.B)  { benchmarkTest(16, b) }
func BenchmarkTest32(b *testing.B)  { benchmarkTest(32, b) }
func BenchmarkTest64(b *testing.B)  { benchmarkTest(64, b) }
func BenchmarkTest128(b *testing.B) { benchmarkTest(128, b) }

func BenchmarkTestRec2(b *testing.B)   { benchmarkTestRec(2, b) }
func BenchmarkTestRec4(b *testing.B)   { benchmarkTestRec(4, b) }
func BenchmarkTestRec8(b *testing.B)   { benchmarkTestRec(8, b) }
func BenchmarkTestRec16(b *testing.B)  { benchmarkTestRec(16, b) }
func BenchmarkTestRec32(b *testing.B)  { benchmarkTestRec(32, b) }
func BenchmarkTestRec64(b *testing.B)  { benchmarkTestRec(64, b) }
func BenchmarkTestRec128(b *testing.B) { benchmarkTestRec(128, b) }

func benchmarkTestRec(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		mysteriousFunctionRec(n)
	}
}

func benchmarkTest(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		mysteriousFunction(n)
	}
}

func Test_mysteriousFunctionRec(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{name: "-1", args: args{n: -1}, want: big.NewInt(1)},
		{name: "0", args: args{n: 0}, want: big.NewInt(1)},
		{name: "1", args: args{n: 1}, want: big.NewInt(1)},
		{name: "2", args: args{n: 2}, want: big.NewInt(2)},
		{name: "3", args: args{n: 3}, want: big.NewInt(6)},
		{name: "4", args: args{n: 4}, want: big.NewInt(24)},
		{name: "70", args: args{n: 70}, want: fromString("11978571669969891796072783721689098736458938142546425857555362864628009582789845319680000000000000000")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mysteriousFunctionRec(tt.args.n); got.Cmp(tt.want) != 0 {
				t.Errorf("mysteriousFunctionRec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysteriousFunction(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{name: "-1", args: args{n: -1}, want: big.NewInt(1)},
		{name: "0", args: args{n: 0}, want: big.NewInt(1)},
		{name: "1", args: args{n: 1}, want: big.NewInt(1)},
		{name: "2", args: args{n: 2}, want: big.NewInt(2)},
		{name: "3", args: args{n: 3}, want: big.NewInt(6)},
		{name: "4", args: args{n: 4}, want: big.NewInt(24)},
		{name: "70", args: args{n: 70}, want: fromString("11978571669969891796072783721689098736458938142546425857555362864628009582789845319680000000000000000")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mysteriousFunction(tt.args.n); got.Cmp(tt.want) != 0 {
				t.Errorf("mysteriousFunctionRec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func fromString(s string) *big.Int {
	i, _ := big.NewInt(1).SetString(s, 10)
	return i
}
