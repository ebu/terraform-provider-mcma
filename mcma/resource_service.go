package mcma

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"time"

	mcmamodel "github.com/ebu/mcma-libraries-go/model"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		Description: "Service data registered in an MCMA Service Registry",

		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "The MCMA type of resource. This value will always be 'Service'.",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the service. MCMA IDs are always absolute urls.",
				Computed:    true,
			},
			"date_created": {
				Type:        schema.TypeString,
				Description: "The date and time at which the service data was created.",
				Computed:    true,
			},
			"date_modified": {
				Type:        schema.TypeString,
				Description: "The date and time at which the service data was last modified.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the service",
				Required:    true,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Description: "The type of authentication the service uses, e.g. AWS4",
				Optional:    true,
			},
			"job_type": {
				Type:        schema.TypeString,
				Description: "The type of job the service processes, if any. Most MCMA services will handle some kind of job, but not all of them have to.",
				Optional:    true,
			},
			"resource": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "The MCMA type of resource. This value will always be 'ResourceEndpoint'.",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the resource endpoint. MCMA IDs are always absolute urls.",
							Computed:    true,
						},
						"date_created": {
							Type:        schema.TypeString,
							Description: "The date and time at which the resource endpoint data was created.",
							Computed:    true,
						},
						"date_modified": {
							Type:        schema.TypeString,
							Description: "The date and time at which the resource endpoint data was last modified.",
							Computed:    true,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Description: "The type of MCMA resource this endpoint handles.",
							Required:    true,
						},
						"http_endpoint": {
							Type:        schema.TypeString,
							Description: "The url for the endpoint.",
							Required:    true,
						},
						"auth_type": {
							Type:        schema.TypeString,
							Description: "The type of authentication expected for this endpoint. This should only be specified if it is different than the auth type specified on the service.",
							Optional:    true,
						},
					},
				},
			},
			"job_profile_ids": {
				Type:        schema.TypeList,
				Description: "The list of IDs for job profiles that can be processed by this service. If the service does not process jobs, this should be empty.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getServiceFromResourceData(d *schema.ResourceData) mcmamodel.Service {
	var resources []mcmamodel.ResourceEndpoint
	for _, r := range d.Get("resource").(*schema.Set).List() {
		resource := r.(map[string]interface{})
		authType := resource["auth_type"].(string)
		resources = append(resources, mcmamodel.ResourceEndpoint{
			Type:         "ResourceEndpoint",
			ResourceType: resource["resource_type"].(string),
			HttpEndpoint: resource["http_endpoint"].(string),
			AuthType:     authType,
		})
	}

	var jobProfileIds []string
	for _, id := range d.Get("job_profile_ids").([]interface{}) {
		jobProfileId := id.(string)
		jobProfileIds = append(jobProfileIds, jobProfileId)
	}

	return mcmamodel.Service{
		Type:          "Service",
		Name:          d.Get("name").(string),
		AuthType:      d.Get("auth_type").(string),
		JobType:       d.Get("job_type").(string),
		JobProfileIds: jobProfileIds,
		Resources:     resources,
	}
}

func resourceServiceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}

	serviceId := d.Id()
	resource, err := resourceManager.Get(reflect.TypeOf(mcmamodel.Service{}), serviceId)
	if err != nil {
		return diag.Errorf("error getting service with id %s: %s", serviceId, err)
	}
	if resource == nil {
		return diag.Diagnostics{}
	}

	service := resource.(mcmamodel.Service)

	_ = d.Set("type", service.Type)
	_ = d.Set("id", service.Id)
	_ = d.Set("date_created", service.DateCreated.Format(time.RFC3339))
	_ = d.Set("date_modified", service.DateModified.Format(time.RFC3339))
	_ = d.Set("name", service.Name)
	_ = d.Set("auth_type", service.AuthType)
	_ = d.Set("job_profile_ids", service.JobProfileIds)

	var resources []map[string]interface{}
	for _, resourceEndpoint := range service.Resources {
		r := make(map[string]interface{})
		r["type"] = "ResourceEndpoint"
		r["resource_type"] = resourceEndpoint.ResourceType
		r["http_endpoint"] = resourceEndpoint.HttpEndpoint
		r["auth_type"] = resourceEndpoint.AuthType
		resources = append(resources, r)
	}
	if err = d.Set("resource", resources); err != nil {
		return diag.Errorf("error setting resources for service with id %s: %s", serviceId, err)
	}

	return diag.Diagnostics{}
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}

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
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}

	service := getServiceFromResourceData(d)
	service.Id = d.Id()

	_, err := resourceManager.Update(service)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceRead(ctx, d, m)
}

func resourceServiceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}

	err := resourceManager.Delete(reflect.TypeOf(mcmamodel.Service{}), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
