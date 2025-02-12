package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Common models used across multiple data sources
type searchModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
