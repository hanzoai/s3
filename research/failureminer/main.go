package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"time"
)

type Telemetry struct {
	Timestamp        time.Time `json:"timestamp"`
	Device           string    `json:"device"`
	CriticalWarning  int       `json:"critical_warning"`
	Temperature      int       `json:"temperature"`
	AvailableSpare   int       `json:"available_spare"`
	MediaErrors      int       `json:"media_errors"`
	NumErrLogEntries int       `json:"num_err_log_entries"`
}

const (
	MaxMediaErrors       = 10
	MinAvailableSparePct = 10
)

var validNVMeDeviceRe = regexp.MustCompile(`^/dev/nvme\d+(n\d+)?$`)

func validateDevicePath(device string) error {
	if !validNVMeDeviceRe.MatchString(device) {
		return errors.New("invalid NVMe device path: expected /dev/nvme<N> or /dev/nvme<N>n<M>")
	}
	return nil
}

func parseSmartLogJSON(raw []byte, device string) (Telemetry, error) {
	var t Telemetry
	if err := json.Unmarshal(raw, &t); err != nil {
		return Telemetry{}, fmt.Errorf("failed to parse smart-log json: %w", err)
	}
	t.Timestamp = time.Now().UTC()
	t.Device = device
	return t, nil
}

func isFailureImminent(t Telemetry) bool {
	// Note: AvailableSpare is a percentage (0-100). 
	// Failure is imminent if spare is below the threshold.
	return t.CriticalWarning > 0 || t.MediaErrors > MaxMediaErrors || (t.AvailableSpare < MinAvailableSparePct)
}

func collectTelemetry(device string) (Telemetry, error) {
	if err := validateDevicePath(device); err != nil {
		return Telemetry{}, err
	}
	// Read-only invocation; no --privileged required when user has read group access to device node
	cmd := exec.Command("nvme", "smart-log", device, "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return Telemetry{}, err
	}
	return parseSmartLogJSON(out, device)
}

func main() {
	device := flag.String("device", "/dev/nvme0n1", "Read-only NVMe device path")
	flag.Parse()

	data, err := collectTelemetry(*device)
	if err != nil {
		log.Fatalf("Failed to collect telemetry: %v", err)
	}

	if isFailureImminent(data) {
		fmt.Printf("[ALERT] SSD Failure Imminent on %s (critical_warning=%d, media_errors=%d, spare=%d%%)\n",
			*device, data.CriticalWarning, data.MediaErrors, data.AvailableSpare)
	}

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal output: %v", err)
	}
	fmt.Println(string(output))
}
