package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// InitializeFileLogger sets up logging.
func InitializeFileLogger(ctx context.Context) (context.Context, error) {
	// Set up context with subsystem.
	ctx = tflog.NewSubsystem(ctx, "centreon")

	// Add provider metadata to all log entries.
	ctx = tflog.SetField(ctx, "provider", "centreon")

	return ctx, nil
}
