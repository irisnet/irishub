package types

import "testing"

func TestValidateTags(t *testing.T) {
	type args struct {
		tags []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "right case", wantErr: false, args: args{tags: []string{"service1"}}},
		{name: "tag contains chinese", wantErr: false, args: args{tags: []string{"服务1"}}},
		{name: "tag contains underscore", wantErr: false, args: args{tags: []string{"service_1"}}},
		{name: "tag contains special characters", wantErr: false, args: args{tags: []string{"service$1"}}},
		{name: "tag is empty ", wantErr: true, args: args{tags: []string{""}}},
		{name: "tag is blank", wantErr: true, args: args{tags: []string{" "}}},
		{name: "tag contains space", wantErr: true, args: args{tags: []string{"service 1"}}},
		{name: "tag is too long", wantErr: true, args: args{tags: []string{"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTags(tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("ValidateTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
