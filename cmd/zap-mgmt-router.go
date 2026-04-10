// Copyright (c) 2026 Hanzo AI, Inc.
//
// This file is part of Hanzo S3 Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/minio/mux"
)

// registerManagementRouter registers /v1/s3/* management HTTP routes.
func registerManagementRouter(router *mux.Router) {
	mgmt := router.PathPrefix("/v1/s3").Subrouter()

	mgmt.Methods(http.MethodGet).Path("/buckets").HandlerFunc(httpTraceAll(mgmtListBucketsHandler))
	mgmt.Methods(http.MethodGet).Path("/health").HandlerFunc(httpTraceAll(mgmtHealthHandler))
	mgmt.Methods(http.MethodGet).Path("/status").HandlerFunc(httpTraceAll(mgmtStatusHandler))
}

// mgmtListBucketsHandler returns a JSON list of buckets.
func mgmtListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		writeErrorJSON(w, http.StatusServiceUnavailable, "object layer not ready")
		return
	}

	buckets, err := objAPI.ListBuckets(r.Context(), BucketOptions{})
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	type bucketEntry struct {
		Name    string    `json:"name"`
		Created time.Time `json:"created"`
	}

	out := make([]bucketEntry, len(buckets))
	for i, b := range buckets {
		out[i] = bucketEntry{Name: b.Name, Created: b.Created}
	}

	writeJSON(w, http.StatusOK, out)
}

// mgmtHealthHandler returns health status.
func mgmtHealthHandler(w http.ResponseWriter, r *http.Request) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		writeErrorJSON(w, http.StatusServiceUnavailable, "not ready")
		return
	}

	result := objAPI.Health(r.Context(), HealthOptions{})
	status := http.StatusOK
	if !result.Healthy {
		status = http.StatusServiceUnavailable
	}

	writeJSON(w, status, map[string]any{
		"healthy":      result.Healthy,
		"healthy_read": result.HealthyRead,
	})
}

// mgmtStatusHandler returns server info.
func mgmtStatusHandler(w http.ResponseWriter, r *http.Request) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		writeErrorJSON(w, http.StatusServiceUnavailable, "not ready")
		return
	}

	info := objAPI.StorageInfo(r.Context(), false)

	writeJSON(w, http.StatusOK, map[string]any{
		"version":  ReleaseTag,
		"runtime":  runtime.Version(),
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"disks":    len(info.Disks),
		"uptime":   time.Since(globalBootTime).String(),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErrorJSON(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
