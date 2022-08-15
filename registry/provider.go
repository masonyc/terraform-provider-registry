package registry

import (
	"context"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"azure_tenant_id": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"azure_client_id": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"azure_client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"registry_base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"registry_resources": resourceRegistryResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"registry_resources": dataSourceResources(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client_secret := d.Get("azure_client_secret").(string)
	client_id := d.Get("azure_client_id").(string)
	tenant_id := d.Get("azure_tenant_id").(string)
	cred, err := confidential.NewCredFromSecret(client_secret)
	if err != nil {
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

	return accessToken.AccessToken, nil
}
