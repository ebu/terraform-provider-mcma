package mcma

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	mcma "github.com/ebu/mcma-libraries-go/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MCMA_SERVICES_URL", nil),
			},
			"authType": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MCMA_SERVICES_AUTH_TYPE", nil),
			},
			"authContext": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MCMA_SERVICES_AUTH_CONTEXT", nil),
			},
		},
		ResourcesMap:   map[string]*schema.Resource{
			"service": resourceService(),
			"job_profile": resourceJobProfile(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	authType := d.Get("auth_type").(string)
	authContext := d.Get("auth_context").(string)

	if url == "" {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Unable to create MCMA ResourceManager",
				Detail:        "Url for MCMA service registry was not provided, and the MCMA_SERVICES_URL environment variable is not set.",
			},
		}
	}

	var resourceManager mcma.ResourceManager
	if authType != "" && authContext != "" {
		resourceManager = mcma.NewResourceManager(url, authType, authContext)
	} else {
		resourceManager = mcma.NewResourceManagerNoAuth(url)
	}

	return resourceManager, nil
}
