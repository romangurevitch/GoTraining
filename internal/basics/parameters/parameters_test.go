package parameters

import (
	"testing"
)

func BenchmarkPointerParameter_1(b *testing.B) { benchPointer(b, 1) }
func BenchmarkValueParameter_1(b *testing.B)   { benchValue(b, 1) }

func BenchmarkPointerParameter_100(b *testing.B) { benchPointer(b, 100) }
func BenchmarkValueParameter_100(b *testing.B)   { benchValue(b, 100) }

func BenchmarkPointerParameter_10000(b *testing.B) { benchPointer(b, 10000) }
func BenchmarkValueParameter_10000(b *testing.B)   { benchValue(b, 10000) }

func BenchmarkPointerParameter_1000000(b *testing.B) { benchPointer(b, 1000000) }
func BenchmarkValueParameter_1000000(b *testing.B)   { benchValue(b, 10000000) }

func BenchmarkPointerParameter_100000000(b *testing.B) { benchPointer(b, 100000000) }
func BenchmarkValueParameter_100000000(b *testing.B)   { benchValue(b, 100000000) }

func benchPointer(b *testing.B, n int) {
	// nolint: lll
	l := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque quam risus, tincidunt pretium diam ac, cursus euismod risus. Donec scelerisque turpis nunc, eu ullamcorper metus mollis ut. Proin vulputate vehicula urna ac facilisis. Curabitur mi nunc, dapibus eu ipsum vitae, ornare malesuada quam. Fusce tempus tincidunt nulla, vel finibus nisi eleifend nec. Duis finibus elit eu tellus hendrerit eleifend. Duis vestibulum velit id dolor tempor fringilla. Pellentesque pulvinar, urna quis mollis cursus, diam dui interdum lacus, sit amet tincidunt massa nibh id urna. Morbi dui felis, gravida a sapien id, posuere hendrerit leo. Cras ac viverra velit. Mauris eu finibus nibh, at pretium lorem. Donec a condimentum velit. Aenean lobortis gravida ligula. Vestibulum placerat feugiat magna ut porttitor. Nullam eget purus laoreet, malesuada mauris ac, mollis ligula.`
	t := Currency{
		Code:  l,
		Money: 9999999999999999.9999,
		Fx:    1234567890987654321,
	}
	for i := 0; i < b.N; i++ {
		pointerParameter(t, n)
	}
}

func benchValue(b *testing.B, n int) {
	// nolint: lll
	l := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque quam risus, tincidunt pretium diam ac, cursus euismod risus. Donec scelerisque turpis nunc, eu ullamcorper metus mollis ut. Proin vulputate vehicula urna ac facilisis. Curabitur mi nunc, dapibus eu ipsum vitae, ornare malesuada quam. Fusce tempus tincidunt nulla, vel finibus nisi eleifend nec. Duis finibus elit eu tellus hendrerit eleifend. Duis vestibulum velit id dolor tempor fringilla. Pellentesque pulvinar, urna quis mollis cursus, diam dui interdum lacus, sit amet tincidunt massa nibh id urna. Morbi dui felis, gravida a sapien id, posuere hendrerit leo. Cras ac viverra velit. Mauris eu finibus nibh, at pretium lorem. Donec a condimentum velit. Aenean lobortis gravida ligula. Vestibulum placerat feugiat magna ut porttitor. Nullam eget purus laoreet, malesuada mauris ac, mollis ligula.`
	t := Currency{
		Code:  l,
		Money: 9999999999999999.9999,
		Fx:    1234567890987654321,
	}
	for i := 0; i < b.N; i++ {
		valueParameter(t, n)
	}
}
