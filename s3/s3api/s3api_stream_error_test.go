package s3api

import (
	"context"
	"fmt"
	"testing"
)

func TestShouldWriteStreamingErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "wrapped context canceled",
			err:      &StreamError{Err: context.Canceled},
			expected: false,
		},
		{
			name:     "grpc canceled",
			err:      fmt.Errorf("Canceled: client connection is closing"),
			expected: false,
		},
		{
			name:     "wrapped grpc canceled",
			err:      &StreamError{Err: fmt.Errorf("Canceled: client connection is closing")},
			expected: false,
		},
		{
			name:     "deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: true,
		},
		{
			name:     "wrapped deadline exceeded",
			err:      &StreamError{Err: context.DeadlineExceeded},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldWriteStreamingErrorResponse(tt.err); got != tt.expected {
				t.Fatalf("shouldWriteStreamingErrorResponse(%v) = %v, want %v", tt.err, got, tt.expected)
			}
		})
	}
}
