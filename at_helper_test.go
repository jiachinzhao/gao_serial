package gao_serial

import (
	"os"
	"testing"
	"time"
)

func TestGetIMEI(t *testing.T) {
	port := os.Getenv("PORT")
	t.Logf("port: %s", port)
	type args struct {
		gs *GaoSerial
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{NewGaoSerial(5 * time.Second)}, want: "861164036724633", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.gs.Open(port, 115200); err != nil {
				t.Errorf("open %v", err)
				return
			}
			got, err := GetIMEI(tt.args.gs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIMEI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIMEI() = %v, want %v", got, tt.want)
			}
		})
	}
}
