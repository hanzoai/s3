package failureminer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// LoadSMARTCSV loads external NVMe/SMART telemetry from a CSV file compatible
// with standard Backblaze and Alibaba NVMe SMART export schemas.
//
// Recognized column aliases (case-insensitive):
//   - drive_id / serial_number / disk_id
//   - timestamp_unix / timestamp / date (RFC3339, YYYY-MM-DD, YYYY-MM-DD HH:MM:SS, or Unix seconds)
//   - tenant_id (optional; defaults to "external-default")
//   - vendor (optional; inferred from model prefix or "external")
//   - model (optional; defaults to "unknown_model")
//   - failure (0/1 or false/true)
//   - failure_mode (optional)
//   - media_errors / nvme_media_errors / smart_187_raw
//   - num_err_log_entries / nvme_error_log_entries / smart_197_raw
//   - available_spare_pct / available_spare / smart_232_normalized
//   - percentage_used / smart_233_raw / smart_177_raw
//   - critical_warning_mask / critical_warning
//   - temperature_c / temperature / smart_194_raw
//   - unsafe_shutdowns / smart_192_raw
//   - controller_busy_time
//   - io_timeout_count / smart_188_raw
//   - p99_read_latency_us / read_latency_p99_us
//   - p99_write_latency_us / write_latency_p99_us
//   - crc_errors / smart_199_raw
//   - reallocated_sectors / smart_5_raw
//   - power_on_hours / smart_9_raw
func LoadSMARTCSV(path string) (FleetDataset, error) {
	f, err := os.Open(path)
	if err != nil {
		return FleetDataset{}, fmt.Errorf("open smart csv: %w", err)
	}
	defer f.Close()
	return ParseSMARTCSV(f, path)
}

// ParseSMARTCSV parses external NVMe/SMART CSV telemetry from an io.Reader.
func ParseSMARTCSV(r io.Reader, name string) (FleetDataset, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true

	header, err := cr.Read()
	if err != nil {
		return FleetDataset{}, fmt.Errorf("read csv header: %w", err)
	}
	colIdx := make(map[string]int, len(header))
	for i, h := range header {
		colIdx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	findCol := func(aliases ...string) int {
		for _, a := range aliases {
			if idx, ok := colIdx[a]; ok {
				return idx
			}
		}
		return -1
	}

	idxDrive := findCol("drive_id", "serial_number", "disk_id")
	idxTS := findCol("timestamp_unix", "timestamp", "date", "dt")
	if idxDrive < 0 || idxTS < 0 {
		return FleetDataset{}, fmt.Errorf("csv missing required drive_id/serial_number or timestamp/date column")
	}

	idxTenant := findCol("tenant_id")
	idxVendor := findCol("vendor")
	idxModel := findCol("model")
	idxFailure := findCol("failure", "failed", "label")
	idxMode := findCol("failure_mode")
	idxMedia := findCol("media_errors", "nvme_media_errors", "smart_187_raw")
	idxErrLog := findCol("num_err_log_entries", "nvme_error_log_entries", "smart_197_raw")
	idxSpare := findCol("available_spare_pct", "available_spare", "smart_232_normalized")
	idxPctUsed := findCol("percentage_used", "smart_233_raw", "smart_177_raw")
	idxCritWarn := findCol("critical_warning_mask", "critical_warning")
	idxTemp := findCol("temperature_c", "temperature", "smart_194_raw")
	idxUnsafe := findCol("unsafe_shutdowns", "smart_192_raw")
	idxBusy := findCol("controller_busy_time")
	idxTimeout := findCol("io_timeout_count", "smart_188_raw")
	idxReadLat := findCol("p99_read_latency_us", "read_latency_p99_us")
	idxWriteLat := findCol("p99_write_latency_us", "write_latency_p99_us")
	idxCRC := findCol("crc_errors", "smart_199_raw")
	idxRealloc := findCol("reallocated_sectors", "smart_5_raw")
	idxPOH := findCol("power_on_hours", "smart_9_raw")

	type rowEntry struct {
		sample  DriveTelemetrySample
		failed  bool
		modeStr string
	}
	byDrive := make(map[string][]rowEntry)

	getStr := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}
	getFloat := func(row []string, idx int, def float64) float64 {
		s := getStr(row, idx)
		if s == "" {
			return def
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return def
		}
		return v
	}
	getUint := func(row []string, idx int) uint64 {
		s := getStr(row, idx)
		if s == "" {
			return 0
		}
		if v, err := strconv.ParseUint(s, 10, 64); err == nil {
			return v
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil && f >= 0 {
			return uint64(f)
		}
		return 0
	}

	lineNum := 1
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return FleetDataset{}, fmt.Errorf("read csv line %d: %w", lineNum+1, err)
		}
		lineNum++

		driveID := getStr(row, idxDrive)
		if driveID == "" {
			continue
		}
		ts, err := parseTimestamp(getStr(row, idxTS))
		if err != nil {
			return FleetDataset{}, fmt.Errorf("csv line %d invalid timestamp: %w", lineNum, err)
		}

		model := getStr(row, idxModel)
		if model == "" {
			model = "external_nvme"
		}
		vendor := getStr(row, idxVendor)
		if vendor == "" {
			if parts := strings.Fields(model); len(parts) > 0 {
				vendor = strings.ToLower(parts[0])
			} else {
				vendor = "external"
			}
		}
		tenant := getStr(row, idxTenant)
		if tenant == "" {
			tenant = "external-default"
		}

		failVal := strings.ToLower(getStr(row, idxFailure))
		isFail := failVal == "1" || failVal == "true" || failVal == "yes"

		spare := getFloat(row, idxSpare, 100.0)
		if spare > 100.0 {
			spare = 100.0
		}
		if spare < 0.0 {
			spare = 0.0
		}

		sample := DriveTelemetrySample{
			TimestampUnix:       ts,
			DriveID:             driveID,
			TenantID:            tenant,
			Vendor:              vendor,
			Model:               model,
			MediaErrors:         getUint(row, idxMedia),
			NumErrLogEntries:    getUint(row, idxErrLog),
			AvailableSparePct:   spare,
			PercentageUsed:      getFloat(row, idxPctUsed, 0.0),
			CriticalWarningMask: uint8(getUint(row, idxCritWarn)),
			TemperatureC:        getFloat(row, idxTemp, 40.0),
			UnsafeShutdowns:     getUint(row, idxUnsafe),
			ControllerBusyTime:  getUint(row, idxBusy),
			IOTimeoutCount:      getUint(row, idxTimeout),
			P99ReadLatencyUs:    getFloat(row, idxReadLat, 300.0),
			P99WriteLatencyUs:   getFloat(row, idxWriteLat, 450.0),
			CRCErrors:           getUint(row, idxCRC),
			ReallocatedSectors:  getUint(row, idxRealloc),
			PowerOnHours:        getUint(row, idxPOH),
		}

		byDrive[driveID] = append(byDrive[driveID], rowEntry{
			sample:  sample,
			failed:  isFail,
			modeStr: getStr(row, idxMode),
		})
	}

	driveIDs := make([]string, 0, len(byDrive))
	for id := range byDrive {
		driveIDs = append(driveIDs, id)
	}
	sort.Strings(driveIDs)

	drives := make([]DriveSeries, 0, len(driveIDs))
	var stepHours int = 24
	for _, id := range driveIDs {
		entries := byDrive[id]
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].sample.TimestampUnix < entries[j].sample.TimestampUnix
		})
		if len(entries) == 0 {
			continue
		}
		if len(entries) >= 2 {
			dt := entries[1].sample.TimestampUnix - entries[0].sample.TimestampUnix
			if dt >= 3600 {
				stepHours = int(dt / 3600)
			}
		}

		var (
			failed      bool
			failTS      int64
			mode        = FailureModeNone
			samples     = make([]DriveTelemetrySample, 0, len(entries))
		)
		for _, e := range entries {
			if e.failed && !failed {
				failed = true
				failTS = e.sample.TimestampUnix
				if e.modeStr != "" {
					mode = FailureMode(e.modeStr)
				} else {
					mode = FailureModeMediaErrorRamp
				}
				// Stop collecting after failure timestamp to preserve causal pre-failure prefix.
				break
			}
			samples = append(samples, e.sample)
		}
		if len(samples) == 0 {
			continue
		}
		meta := entries[0].sample.Meta()
		drives = append(drives, DriveSeries{
			Meta:            meta,
			Failed:          failed,
			FailureMode:     mode,
			FailureTimeUnix: failTS,
			Samples:         samples,
		})
	}

	if name == "" {
		name = "external_smart_csv"
	}
	return FleetDataset{
		Name:      name,
		StepHours: stepHours,
		Drives:    drives,
	}, nil
}

// LoadSMARTJSON loads a FleetDataset previously serialized as JSON.
func LoadSMARTJSON(path string) (FleetDataset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FleetDataset{}, fmt.Errorf("read smart json: %w", err)
	}
	var ds FleetDataset
	if err := json.Unmarshal(data, &ds); err != nil {
		return FleetDataset{}, fmt.Errorf("decode smart json: %w", err)
	}
	for i := range ds.Drives {
		sort.Slice(ds.Drives[i].Samples, func(a, b int) bool {
			return ds.Drives[i].Samples[a].TimestampUnix < ds.Drives[i].Samples[b].TimestampUnix
		})
	}
	return ds, nil
}

func parseTimestamp(raw string) (int64, error) {
	if raw == "" {
		return 0, fmt.Errorf("empty timestamp")
	}
	if sec, err := strconv.ParseInt(raw, 10, 64); err == nil {
		if sec > 1e12 {
			// Milliseconds -> seconds.
			return sec / 1000, nil
		}
		return sec, nil
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC().Unix(), nil
		}
	}
	return 0, fmt.Errorf("unrecognized timestamp format %q", raw)
}
