package mcma

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	mcmaclient "github.com/ebu/mcma-libraries-go/client"
)

func resourceMcmaResource() *schema.Resource {
	return &schema.Resource{
		Description: "An arbitrary MCMA resource managed through a REST API",

		CreateContext: resourceMcmaResourceCreate,
		ReadContext:   resourceMcmaResourceRead,
		UpdateContext: resourceMcmaResourceUpdate,
		DeleteContext: resourceMcmaResourceDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "The MCMA type of resource.",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the service. MCMA IDs are always absolute urls.",
				Computed:    true,
			},
			"resource_json": {
				Type:        schema.TypeString,
				Description: "The JSON of the object to be created",
				Required:    true,
			},
		},
	}
}

func getMcmaResourceFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	resourceJson := d.Get("resource_json").(string)

	resourceMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(resourceJson), &resourceMap)
	if err != nil {
		return nil, fmt.Errorf("error parsing json to map: %v", err)
	}

	resourceMap["id"] = d.Id()
	resourceMap["@type"] = d.Get("type").(string)

	return resourceMap, nil
}

func resourceMcmaResourceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	resourceType := d.Get("type").(string)
	resourceId := d.Id()
	resource, err := resourceManager.GetResource(resourceType, resourceId)
	if err != nil {
		return diag.Errorf("error getting resource of type %s with id %s: %s", resourceType, resourceId, err)
	}
	if resource == nil {
		return diag.Errorf("resource with type %s and id %s not found", resourceType, resourceId)
	}

	_ = d.Set("type", resource["@type"])
	delete(resource, "@type")
	_ = d.Set("id", resource["id"])
	delete(resource, "id")

	delete(resource, "dateCreated")
	delete(resource, "dateModified")

	jsonBytes, err := json.Marshal(resource)
	if err != nil {
		return diag.Errorf("error parsing json for resource of type %s with id %s: %s", resourceType, resourceId, err)
	}
	_ = d.Set("resource_json", string(jsonBytes))

	return diag.Diagnostics{}
}

func resourceMcmaResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	resource, err := getMcmaResourceFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createdResource, err := resourceManager.Create(resource)
	if err != nil {
		return diag.FromErr(err)
	}
	createdResourceMap := createdResource.(map[string]interface{})

	d.SetId(createdResourceMap["id"].(string))

	return resourceMcmaResourceRead(ctx, d, m)
}

func resourceMcmaResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	resource, err := getMcmaResourceFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = resourceManager.Update(resource)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMcmaResourceRead(ctx, d, m)
}

func resourceMcmaResourceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	err := resourceManager.DeleteResource(d.Get("type").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
