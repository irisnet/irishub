package types

import "testing"

func TestValidatePoolName(t *testing.T) {
	type args struct {
		poolName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				poolName: "1pool",
			},
			wantErr: true,
		},
		{
			name: "test case 2",
			args: args{
				poolName: "pool&",
			},
			wantErr: true,
		},
		{
			name: "test case 3",
			args: args{
				poolName: "12345678901234567890123456789012345678901234567890123456789012345678901234567890",
			},
			wantErr: true,
		},
		{
			name: "test case 4",
			args: args{
				poolName: "BUSD-IRIS",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePoolName(tt.args.poolName); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePoolName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
