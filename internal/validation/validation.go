package validation

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SNMPVersionValidator validates that SNMP version is one of: 1, 2c, or 3.
type SNMPVersionValidator struct{}

func (v SNMPVersionValidator) Description(ctx context.Context) string {
	return "SNMP version must be one of: 1, 2c, or 3"
}

func (v SNMPVersionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v SNMPVersionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value != "1" && value != "2c" && value != "3" {
		resp.Diagnostics.AddError(
			"Invalid SNMP Version",
			fmt.Sprintf("SNMP version must be one of: 1, 2c, or 3, got: %s", value),
		)
	}
}

// HostnameOrIPValidator validates that a string is either a valid hostname or IP address.
type HostnameOrIPValidator struct{}

func (v HostnameOrIPValidator) Description(ctx context.Context) string {
	return "value must be a valid hostname or IP address"
}

func (v HostnameOrIPValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v HostnameOrIPValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Check if it's an IP address.
	if ip := net.ParseIP(value); ip != nil {
		return
	}

	// Check if it's a valid hostname.
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !hostnameRegex.MatchString(value) {
		resp.Diagnostics.AddError(
			"Invalid Address",
			fmt.Sprintf("%s is neither a valid IP address nor a valid hostname", value),
		)
	}
}

// GeoCoordsValidator validates geographic coordinates in format "lat,long".
type GeoCoordsValidator struct{}

func (v GeoCoordsValidator) Description(ctx context.Context) string {
	return "value must be valid geographic coordinates in format 'latitude,longitude'"
}

func (v GeoCoordsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v GeoCoordsValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Geographic Coordinates",
			fmt.Sprintf("Coordinates must be in format 'latitude,longitude', got: %s", value),
		)
		return
	}

	// Simple validation of lat/long ranges
	lat := strings.TrimSpace(parts[0])
	long := strings.TrimSpace(parts[1])
	latRegex := regexp.MustCompile(`^-?([0-9]|[1-8][0-9]|90)(\.[0-9]+)?$`)
	longRegex := regexp.MustCompile(`^-?([0-9]|[1-9][0-9]|1[0-7][0-9]|180)(\.[0-9]+)?$`)

	if !latRegex.MatchString(lat) {
		resp.Diagnostics.AddError(
			"Invalid Latitude",
			fmt.Sprintf("Latitude must be between -90 and 90 degrees, got: %s", lat),
		)
	}
	if !longRegex.MatchString(long) {
		resp.Diagnostics.AddError(
			"Invalid Longitude",
			fmt.Sprintf("Longitude must be between -180 and 180 degrees, got: %s", long),
		)
	}
}

// NotificationOptionsValidator validates notification options.
type NotificationOptionsValidator struct{}

func (v NotificationOptionsValidator) Description(ctx context.Context) string {
	return "value must be a valid combination of notification options (1=DOWN, 2=UNREACHABLE, 4=RECOVERY, 8=FLAPPING, 16=DOWNTIME_SCHEDULED)"
}

func (v NotificationOptionsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v NotificationOptionsValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueInt64()
	validBits := int64(1 | 2 | 4 | 8 | 16) // All valid options combined.

	if value < 0 || value > validBits {
		resp.Diagnostics.AddError(
			"Invalid Notification Options",
			fmt.Sprintf("Notification options must be a valid combination of: 1=DOWN, 2=UNREACHABLE, 4=RECOVERY, 8=FLAPPING, 16=DOWNTIME_SCHEDULED. Got: %d", value),
		)
	}
}
