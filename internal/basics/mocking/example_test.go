package mocking

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/romangurevitch/go-training/internal/basics/mocking/calculator"
	"github.com/romangurevitch/go-training/internal/basics/mocking/calculator/mocks"
)

// TestExampleFunction_Times shows Times(n): the mock asserts the method is
// called exactly n times. If called more or fewer times, the test fails.
func TestExampleFunction_Times(t *testing.T) {
	ctrl := gomock.NewController(t)
	a := mocks.NewMockAdder(ctrl)

	// Expect exactly one call — calling it 0 or 2+ times fails the test.
	a.EXPECT().SingleDigitAdd(1, 2).Return(3, nil).Times(1)

	got, err := ExampleFunction(a, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

// TestExampleFunction_DoAndReturn shows DoAndReturn: execute custom logic
// instead of returning static values. Useful for dynamic responses or
// verifying the arguments passed to the mock.
func TestExampleFunction_DoAndReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	a := mocks.NewMockAdder(ctrl)

	a.EXPECT().
		SingleDigitAdd(gomock.Any(), gomock.Any()).
		DoAndReturn(func(x, y int) (int, error) {
			if x+y >= 10 {
				return 0, errors.New("too high")
			}
			return x + y, nil
		})

	got, err := ExampleFunction(a, 3, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 7 {
		t.Errorf("got %d, want 7", got)
	}
}

// TestExampleFunction_InOrder shows InOrder: enforce that calls happen in a
// specific sequence. If Subtract is called before Add, the test fails.
func TestExampleFunction_InOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	a := mocks.NewMockAdder(ctrl)

	first := a.EXPECT().SingleDigitAdd(1, 2).Return(3, nil)
	a.EXPECT().SingleDigitAdd(3, 4).Return(7, nil).After(first)

	_, _ = ExampleFunction(a, 1, 2)
	_, _ = ExampleFunction(a, 3, 4)
}

func TestExampleFunction(t *testing.T) {
	type args struct {
		adder func() calculator.Adder
		x     int
		y     int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "no error testcase",
			args: args{
				adder: func() calculator.Adder {
					ctrl := gomock.NewController(t)
					a := mocks.NewMockAdder(ctrl)
					a.EXPECT().SingleDigitAdd(1, 2).Return(3, nil)
					return a
				},
				x: 1,
				y: 2,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "error testcase",
			args: args{
				adder: func() calculator.Adder {
					ctrl := gomock.NewController(t)
					a := mocks.NewMockAdder(ctrl)
					a.EXPECT().SingleDigitAdd(gomock.Any(), gomock.Any()).Return(0, errors.New("error"))
					return a
				},
				x: 1,
				y: 2,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleFunction(tt.args.adder(), tt.args.x, tt.args.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExampleFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExampleFunction() got = %v, want %v", got, tt.want)
			}
		})
	}
}
