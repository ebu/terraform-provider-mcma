package mcma

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	mcma "github.com/ebu/mcma-libraries-go/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The url to the services endpoint of the MCMA Service Registry in which you wish to register services and job profiles",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MCMA_SERVICES_URL", nil),
			},
			"auth_type": {
				Type:        schema.TypeString,
				Description: "The authentication type to use for authenticating calls to the MCMA Service Registry",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MCMA_SERVICES_AUTH_TYPE", nil),
			},
			"auth_context": {
				Type:        schema.TypeString,
				Description: "The authentication context to use for authenticating calls to the MCMA Service Registry",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MCMA_SERVICES_AUTH_CONTEXT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mcma_service":     resourceService(),
			"mcma_job_profile": resourceJobProfile(),
			"mcma_resource":    resourceMcmaResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return configure(d)
		},
	}
}

func configure(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	authType := d.Get("auth_type").(string)
	authContext := d.Get("auth_context").(string)

	if url == "" {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create MCMA ResourceManager",
				Detail:   "Url for MCMA service registry was not provided, and the MCMA_SERVICES_URL environment variable is not set.",
			},
		}
	}

	var resourceManager mcma.ResourceManager
	if authType != "" && authContext != "" {
		resourceManager = mcma.NewResourceManager(url, authType, authContext)
	} else {
		resourceManager = mcma.NewResourceManagerNoAuth(url)
	}
	resourceManager.AddAWS4Auth()

	return &resourceManager, nil
}
