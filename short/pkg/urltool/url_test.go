package urltool

import "testing"

func TestGetBasePath(t *testing.T) {

	type args struct {
		Requrl string
	}
	tests := []struct {
		name       string
		args       args
		wantTarget string
		wantErr    bool
	}{
		{
			name:       "normal example",
			args:       args{Requrl: "http://www.kedudu.com/rabbit/"},
			wantTarget: "rabbit",
			wantErr:    false,
		},

		{
			name:       "case 1:",
			args:       args{Requrl: "/abcd/12345"},
			wantTarget: "",
			wantErr:    true,
		},

		{
			name:       "case 2:",
			args:       args{Requrl: "http://www.kedudu.com/rabbit/?age=18&money=1000"},
			wantTarget: "rabbit",
			wantErr:    false,
		},

		{
			name:       "case 3:",
			args:       args{Requrl: ""},
			wantTarget: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTarget, err := GetBasePath(tt.args.Requrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTarget != tt.wantTarget {
				t.Errorf("GetBasePath() = %v, want %v", gotTarget, tt.wantTarget)
			}
		})
	}
}
