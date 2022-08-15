package registry

import (
	"context"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
	"net/http"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	accessToken     string
	configured      bool
	client          *http.Client
	registryBaseUrl string
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"azure_tenant_id": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
			},
			"azure_client_id": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
			},
			"azure_client_secret": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
			},
			"registry_base_url": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	AzureTenantId     types.String `tfsdk:"azure_tenant_id"`
	AzureClientId     types.String `tfsdk:"azure_client_id"`
	AzureClientSecret types.String `tfsdk:"azure_client_secret"`
	RegistryBaseURL   types.String `tfsdk:"registry_base_url"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant_id := config.AzureTenantId.Value
	if tenant_id == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find Azure Tenant Id",
			"Azure Tenant Id cannot be an empty string",
		)
		return
	}
	client_id := config.AzureClientId.Value
	if client_id == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find Azure Client Id",
			"Azure Client Id cannot be an empty string",
		)
		return
	}
	client_secret := config.AzureClientSecret.Value
	if client_secret == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find Azure Client Secret",
			"Azure Client Secret cannot be an empty string",
		)
		return
	}

	registryBaseURL := config.RegistryBaseURL.Value

	if registryBaseURL == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find registry base URL",
			"Registry base URL key cannot be an empty string",
		)
		return
	}
	cred, err := confidential.NewCredFromSecret(client_secret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create a cred from a azure client secret",
			err.Error(),
		)
		return
	}
	confidentialClientApp, err := confidential.New(client_id, cred, confidential.WithAuthority(fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenant_id)))

	scopes := []string{
		"https://serko.co.nz/42f0101f-10b2-453c-87bd-8b3bcabbdc63/.default",
	}
	accessToken, err := confidentialClientApp.AcquireTokenSilent(ctx, scopes)
	if err != nil {
		accessToken, err = confidentialClientApp.AcquireTokenByCredential(context.Background(), scopes)
		if err != nil {
			log.Fatal(err)
		}
	}
	p.client = &http.Client{}
	p.accessToken = accessToken.AccessToken
	p.registryBaseUrl = registryBaseURL
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"registry_resources": resourceRegistryResource{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
