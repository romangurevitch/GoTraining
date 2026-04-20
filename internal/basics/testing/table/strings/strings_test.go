package strings

import "testing"

func TestStrings_AddPrefix(t *testing.T) {
	type fields struct {
		prefix string
	}
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "adds prefix to string",
			fields:  fields{prefix: "hello_"},
			args:    args{str: "world"},
			want:    "hello_world",
			wantErr: false,
		},
		{
			name:    "empty prefix",
			fields:  fields{prefix: ""},
			args:    args{str: "world"},
			want:    "world",
			wantErr: false,
		},
		{
			name:    "empty string returns error",
			fields:  fields{prefix: "hello_"},
			args:    args{str: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "both prefix and string present",
			fields:  fields{prefix: "pre-"},
			args:    args{str: "fix"},
			want:    "pre-fix",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Strings{
				prefix: tt.fields.prefix,
			}
			got, err := s.AddPrefix(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddPrefix() got = %v, want %v", got, tt.want)
			}
		})
	}
}
