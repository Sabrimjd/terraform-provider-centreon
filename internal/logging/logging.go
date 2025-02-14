package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// InitializeFileLogger sets up logging
func InitializeFileLogger(ctx context.Context) (context.Context, error) {
	// Get log file path from TF_LOG_PATH env variable
	logFilePath := os.Getenv("TF_LOG_PATH")
	if logFilePath == "" {
		logFilePath = "terraform-provider-centreon.log"
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return ctx, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Set up context with subsystem
	ctx = tflog.NewSubsystem(ctx, "centreon")

	// Add provider metadata to all log entries
	ctx = tflog.SetField(ctx, "provider", "centreon")

	return ctx, nil
}

// Trace logs a trace message
func Trace(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	for _, fields := range additionalFields {
		for k, v := range fields {
			ctx = tflog.SetField(ctx, k, v)
		}
	}
	tflog.Trace(ctx, msg)
}

// Debug logs a debug message
func Debug(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	for _, fields := range additionalFields {
		for k, v := range fields {
			ctx = tflog.SetField(ctx, k, v)
		}
	}
	tflog.Debug(ctx, msg)
}

// Info logs an info message
func Info(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	for _, fields := range additionalFields {
		for k, v := range fields {
			ctx = tflog.SetField(ctx, k, v)
		}
	}
	tflog.Info(ctx, msg)
}

// Warn logs a warning message
func Warn(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	for _, fields := range additionalFields {
		for k, v := range fields {
			ctx = tflog.SetField(ctx, k, v)
		}
	}
	tflog.Warn(ctx, msg)
}

// Error logs an error message
func Error(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	for _, fields := range additionalFields {
		for k, v := range fields {
			ctx = tflog.SetField(ctx, k, v)
		}
	}
	tflog.Error(ctx, msg)
}
