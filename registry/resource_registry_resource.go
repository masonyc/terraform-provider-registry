package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	registry "terraform-provider-registry/registry/models"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRegistryResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRegistryResourceCreate,
		ReadContext:   resourceRegistryResourceRead,
		UpdateContext: resourceRegistryResourceUpdate,
		DeleteContext: resourceRegistryResourceDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRegistryResourceDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	return nil
}

func resourceRegistryResourceUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	return nil
}

func resourceRegistryResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	access_token := m.(string)
	client := &http.Client{Timeout: 10 * time.Second}
	name := d.Get("name").(string)

	resourceName, err := strconv.Unquote(name)
	requestDTO := registry.ResourceDTO{
		Name: resourceName,
	}

	jsonData, err := json.Marshal(requestDTO)

	requestURL := fmt.Sprintf("%s/resources", "https://api-registry.testing.serko.travel")
	registryRequest, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonData))

	if err != nil {
	}

	registryRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	registryRequest.Header.Set("Content-Type", "application/json")

	response, err := client.Do(registryRequest)
	defer response.Body.Close()

	resource := registry.ResourceDTO{}

	body, err := io.ReadAll(response.Body)
	err = json.Unmarshal(body, &resource)

	d.SetId(resource.Id)
	d.Set("name", resource.Name)

	return nil
}

func resourceRegistryResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
