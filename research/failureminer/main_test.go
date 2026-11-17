package main

import "testing"

func TestValidateDevicePath(t *testing.T) {
	cases := []struct {
		device  string
		wantErr bool
	}{
		{"/dev/nvme0n1", false},
		{"/dev/nvme1", false},
		{"/dev/sda; rm -rf /", true},
		{"../etc/passwd", true},
		{"", true},
	}
	for _, tc := range cases {
		err := validateDevicePath(tc.device)
		if (err != nil) != tc.wantErr {
			t.Fatalf("validateDevicePath(%q) err=%v wantErr=%v", tc.device, err, tc.wantErr)
		}
	}
}

func TestParseSmartLogAndThresholds(t *testing.T) {
	cases := []struct {
		name     string
		payload  string
		imminent bool
		wantErr  bool
	}{
		{
			name:     "healthy",
			payload:  `{"critical_warning":0,"temperature":310,"available_spare":100,"media_errors":0,"num_err_log_entries":0}`,
			imminent: false,
		},
		{
			name:     "degraded_below_threshold",
			payload:  `{"critical_warning":0,"temperature":335,"available_spare":45,"media_errors":5,"num_err_log_entries":2}`,
			imminent: false,
		},
		{
			name:     "failing_media_errors",
			payload:  `{"critical_warning":0,"temperature":345,"available_spare":80,"media_errors":15,"num_err_log_entries":12}`,
			imminent: true,
		},
		{
			name:     "failing_critical_warning",
			payload:  `{"critical_warning":1,"temperature":340,"available_spare":8,"media_errors":2,"num_err_log_entries":4}`,
			imminent: true,
		},
		{
			name:    "invalid_json",
			payload: `{bad-json`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseSmartLogJSON([]byte(tc.payload), "/dev/nvme0n1")
			if (err != nil) != tc.wantErr {
				t.Fatalf("parse error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if isFailureImminent(got) != tc.imminent {
				t.Fatalf("isFailureImminent = %v, want %v (data=%+v)", isFailureImminent(got), tc.imminent, got)
			}
		})
	}
}
