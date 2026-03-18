package entity

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		externalField string
	}
	tests := []struct {
		name string
		args args
		want Entity
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.externalField); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entity_Do(t *testing.T) {
	type fields struct {
		externalField string
		internalField string
	}
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entity{
				ExternalField: tt.fields.externalField,
				internalField: tt.fields.internalField,
			}
			got, err := e.Do(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Do() got = %v, want %v", got, tt.want)
			}
		})
	}
}
