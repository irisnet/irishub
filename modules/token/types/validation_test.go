package types

import "testing"

func TestValidateSymbol(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "right case",
			wantErr: false,
			args:    args{symbol: "btc"},
		},
		{
			name:    "start with a capital letter",
			wantErr: true,
			args:    args{symbol: "Btc"},
		},
		{
			name:    "contain a capital letter",
			wantErr: true,
			args:    args{symbol: "bTc"},
		},
		{
			name:    "less than 3 characters in length",
			wantErr: true,
			args:    args{symbol: "ht"},
		},
		{
			name:    "equal 64 characters in length",
			wantErr: false,
			args:    args{symbol: "btc1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct1"},
		},
		{
			name:    "more than 64 characters in length",
			wantErr: true,
			args:    args{symbol: "btc1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct12"},
		},
		{
			name:    "contain peg",
			wantErr: true,
			args:    args{symbol: "pegeth"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSymbol(tt.args.symbol); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbol() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateKeywords(t *testing.T) {
	type args struct {
		denom string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "right case", args: args{denom: "stake"}, wantErr: false},
		{name: "denom is peg", args: args{denom: "peg"}, wantErr: true},
		{name: "denom is Peg", args: args{denom: "Peg"}, wantErr: false},
		{name: "denom begin with peg", args: args{denom: "pegtoken"}, wantErr: true},
		{name: "denom is ibc", args: args{denom: "ibc"}, wantErr: true},
		{name: "denom is IBC", args: args{denom: "Peg"}, wantErr: false},
		{name: "denom begin with ibc", args: args{denom: "ibctoken"}, wantErr: true},
		{name: "denom is lpt", args: args{denom: "lpt"}, wantErr: true},
		{name: "denom is SWAP", args: args{denom: "SWAP"}, wantErr: false},
		{name: "denom begin with lpt", args: args{denom: "lpttoken"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateKeywords(tt.args.denom); (err != nil) != tt.wantErr {
				t.Errorf("CheckKeywords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
