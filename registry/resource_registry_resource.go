package registry

import (
	registry "PoC.RegistryState/registry/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"net/http"
	"strconv"
)

type resourceRegistryResource struct{}

// Resource Resource schema
func (r resourceRegistryResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

type resourceResource struct {
	p provider
}

// New resource instance
func (r resourceRegistryResource) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceResource{
		p: *(p.(*provider)),
	}, nil
}

// Create a new resource
func (r resourceResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}
	var plan registry.Resource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceName, err := strconv.Unquote(plan.Name.String())
	requestDTO := registry.ResourceDTO{
		Name: resourceName,
	}

	jsonData, err := json.Marshal(requestDTO)

	requestURL := fmt.Sprintf("%s/resources", r.p.registryBaseUrl)
	registryRequest, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonData))

	if err != nil {
		resp.Diagnostics.AddError("HttpClient construct failed", "HttpClient construction failed before sending out the request")
	}

	registryRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.p.accessToken))
	registryRequest.Header.Set("Content-Type", "application/json")

	response, err := r.p.client.Do(registryRequest)
	if err != nil || response.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error creating resource",
			"Could not create resource, unexpected error: "+err.Error(),
		)
		return
	}
	defer response.Body.Close()

	resource := registry.ResourceDTO{}

	body, err := io.ReadAll(response.Body)
	err = json.Unmarshal(body, &resource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Decode Response failed",
			"Could not decode response from GET deployment, unexpected error: "+err.Error(),
		)
		return
	}

	mappedTFResource := registry.Resource{}
	mappedTFResource.Id.Value = resource.Id
	mappedTFResource.Name.Value = resource.Name

	diags = resp.State.Set(ctx, mappedTFResource)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state registry.Resource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := state.Id.Value

	requestURL := fmt.Sprintf("%s/resources/%s", r.p.registryBaseUrl, resourceId)
	registryRequest, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		resp.Diagnostics.AddError("HttpClient construct failed", "HttpClient construction failed before sending out the request")
	}

	registryRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.p.accessToken))
	response, err := r.p.client.Do(registryRequest)
	if err != nil || response.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading Resource",
			"Could not read deployment, unexpected error: "+err.Error(),
		)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	resource := registry.ResourceDTO{}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Decode Response failed",
			"Could not decode response from GET deployment, unexpected error: "+err.Error(),
		)
		return
	}

	mappedTFResource := registry.Resource{}
	mappedTFResource.Id.Value = resource.Id
	mappedTFResource.Name.Value = resource.Name

	// Set state
	diags = resp.State.Set(ctx, mappedTFResource)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan registry.Resource
	var state registry.Resource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = state.Id
	resourceName, _ := strconv.Unquote(plan.Name.String())
	resourceId, _ := strconv.Unquote(plan.Id.String())
	requestDTO := registry.ResourceDTO{
		Id:   resourceId,
		Name: resourceName,
	}

	jsonData, err := json.Marshal(requestDTO)

	requestURL := fmt.Sprintf("%s/resources", r.p.registryBaseUrl)
	registryRequest, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(jsonData))

	if err != nil {
		resp.Diagnostics.AddError("HttpClient construct failed", "HttpClient construction failed before sending out the request")
	}

	registryRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.p.accessToken))
	registryRequest.Header.Set("Content-Type", "application/json")

	response, err := r.p.client.Do(registryRequest)
	if err != nil || response.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not update resource, unexpected error: "+err.Error(),
		)
		return
	}
	defer response.Body.Close()

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state registry.Resource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := state.Id.Value
	requestURL := fmt.Sprintf("%s/resources/%s", r.p.registryBaseUrl, resourceId)
	registryRequest, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	if err != nil {
		resp.Diagnostics.AddError("HttpClient construct failed", "HttpClient construction failed before sending out the request")
	}

	registryRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.p.accessToken))
	registryRequest.Header.Set("Content-Type", "application/json")

	response, err := r.p.client.Do(registryRequest)
	if err != nil || response.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error delete resource",
			"Could not delete resource, unexpected error: "+err.Error(),
		)
		return
	}
	defer response.Body.Close()

	resp.State.RemoveResource(ctx)
}

// Import resource

func (r resourceResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
