package registry

import "github.com/hashicorp/terraform-plugin-framework/types"

type Resource struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type ResourceDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
