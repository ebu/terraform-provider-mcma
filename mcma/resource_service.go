package mcma

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/diag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	mcmaclient "github.com/ebu/mcma-libraries-go/client"
	mcmamodel "github.com/ebu/mcma-libraries-go/model"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,

		Schema: map[string]*schema.Schema{
			"@type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: false,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Required: false,
			},
			"auth_context": {
				Type:     schema.TypeString,
				Required: false,
			},
			"job_type": {
				Type:     schema.TypeString,
				Required: false,
			},
			"resources": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"http_endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auth_type": {
							Type:     schema.TypeString,
							Required: false,
						},
						"auth_context": {
							Type:     schema.TypeString,
							Required: false,
						},
					},
				},
			},
			"job_profile_ids": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"input_locations": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"output_locations": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func getServiceFromResourceData(d *schema.ResourceData) mcmamodel.Service{

	var resources []mcmamodel.ResourceEndpoint
	for _, r := range d.Get("resources").([]interface{}) {
		resource := r.(map[string]interface{})
		resources = append(resources, mcmamodel.ResourceEndpoint{
			Type:         "ResourceEndpoint",
			ResourceType: resource["resource_type"].(string),
			HttpEndpoint: resource["http_endpoint"].(string),
			AuthType:     resource["auth_type"].(string),
			AuthContext:  resource["auth_context"].(string),
		})
	}

	var inputLocations []mcmamodel.Locator
	for _, r := range d.Get("input_locations").([]interface{}) {
		inputLocation := r.(map[string]interface{})
		inputLocations = append(inputLocations, mcmamodel.Locator{
			Type: "Locator",
			Url:  inputLocation["url"].(string),
		})
	}

	var outputLocations []mcmamodel.Locator
	for _, r := range d.Get("resources").([]interface{}) {
		outputLocation := r.(map[string]interface{})
		outputLocations = append(outputLocations, mcmamodel.Locator{
			Type: "Locator",
			Url:  outputLocation["url"].(string),
		})
	}

	return mcmamodel.Service{
		Type:            "Service",
		Name:            d.Get("name").(string),
		AuthType:        d.Get("auth_type").(string),
		AuthContext:     d.Get("auth_context").(string),
		JobType:         d.Get("job_type").(string),
		Resources:       resources,
		InputLocations:  inputLocations,
		OutputLocations: outputLocations,
	}
}

func resourceServiceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	serviceId := d.Id()
	resource, err := resourceManager.Get("Service", serviceId)
	if err != nil {
		return diag.Errorf("error getting service with id %s: %s", serviceId, err)
	}

	service := resource.(mcmamodel.Service)

	_ = d.Set("type", service.Type)
	_ = d.Set("id", service.Id)
	_ = d.Set("date_created", service.DateCreated)
	_ = d.Set("date_modified", service.DateModified)
	_ = d.Set("name", service.Name)
	_ = d.Set("auth_type", service.AuthType)
	_ = d.Set("auth_context", service.AuthContext)
	_ = d.Set("job_profile_ids", service.JobProfileIds)

	var resources []map[string]interface{}
	for _, resource := range service.Resources {
		r := make(map[string]interface{})
		r["type"] = "ResourceEndpoint"
		r["id"] = resource.Id
		r["date_created"] = resource.DateCreated
		r["date_modified"] = resource.DateModified
		r["resource_type"] = resource.ResourceType
		r["http_endpoint"] = resource.HttpEndpoint
		r["auth_type"] = resource.AuthType
		r["auth_context"] = resource.AuthContext
		resources = append(resources, r)
	}
	if err = d.Set("resources", resources); err != nil {
		return diag.Errorf("error setting resources for service with id %s: %s", serviceId, err)
	}

	var inputLocations []map[string]interface{}
	for _, inputLocation := range service.InputLocations {
		l := make(map[string]interface{})
		l["type"] = "Locator"
		l["url"] = inputLocation.Url
		inputLocations = append(inputLocations, l)
	}
	if err = d.Set("input_locations", inputLocations); err != nil {
		return diag.Errorf("error setting input_locations for service with id %s: %s", serviceId, err)
	}

	var outputLocations []map[string]interface{}
	for _, inputLocation := range service.OutputLocations {
		l := make(map[string]interface{})
		l["type"] = "Locator"
		l["url"] = inputLocation.Url
		outputLocations = append(outputLocations, l)
	}
	if err = d.Set("output_locations", outputLocations); err != nil {
		return diag.Errorf("error setting output_locations for service with id %s: %s", serviceId, err)
	}

	return diag.Diagnostics{}
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	service := getServiceFromResourceData(d)
	createdResource, err := resourceManager.Create(service)
	if err != nil {
		return diag.FromErr(err)
	}
	createdService := createdResource.(mcmamodel.Service)

	d.SetId(createdService.Id)

	return resourceServiceRead(ctx, d, m)
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	service := getServiceFromResourceData(d)
	service.Id = d.Id()

	_, err := resourceManager.Update(service)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceRead(ctx, d, m)
}

func resourceServiceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	err := resourceManager.Delete("Service", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
