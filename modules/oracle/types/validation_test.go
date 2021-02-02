package types

import "testing"

func TestValidateFeedName(t *testing.T) {
	type args struct {
		feedName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "contain -",
			args:    args{feedName: "usdt-atom"},
			wantErr: false,
		},
		{
			name:    "contain /",
			args:    args{feedName: "usdt/atom"},
			wantErr: false,
		},
		{
			name:    "contain uppercase letter",
			args:    args{feedName: "USDT-atom"},
			wantErr: false,
		},
		{
			name:    "contain digital",
			args:    args{feedName: "USDT-atom2"},
			wantErr: false,
		},
		{
			name:    "start with a number",
			args:    args{feedName: "2USDT-atom2"},
			wantErr: true,
		},
		{
			name:    "contain special characters",
			args:    args{feedName: "USDT$-atom2"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFeedName(tt.args.feedName); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFeedName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
