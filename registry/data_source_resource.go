package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"net/http"
	registry "terraform-provider-registry/registry/models"
	"time"
)

func dataSourceResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcesRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceResourcesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	access_token := m.(string)
	client := &http.Client{Timeout: 10 * time.Second}
	// Get current state
	resourceId := d.Get("id").(int)

	requestURL := fmt.Sprintf("%s/resources/%s", "https://api-registry.testing.serko.travel", resourceId)
	registryRequest, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
	}

	registryRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	response, err := client.Do(registryRequest)
	if err != nil || response.StatusCode != 200 {
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	resource := registry.ResourceDTO{}
	err = json.Unmarshal(body, &resource)
	if err != nil {
	}

	// Set state
	// always run
	d.Set("name", resource.Name)
	d.SetId(resource.Id)

	return nil
}
