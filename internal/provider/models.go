package provider

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Common models used across multiple data sources.
type searchModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// buildSearchQuery converts a searchModel into the JSON search object
// expected by the Centreon API ({"<field>":"<value>"}). Values are JSON-
// marshalled so quotes or special characters in the value cannot break the
// query. Returns "" when the search is not fully specified.
func buildSearchQuery(s *searchModel) string {
	if s == nil || s.Name.IsNull() || s.Value.IsNull() {
		return ""
	}
	payload, err := json.Marshal(map[string]string{s.Name.ValueString(): s.Value.ValueString()})
	if err != nil {
		return ""
	}
	return string(payload)
}
