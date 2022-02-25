package mcma

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	mcma "github.com/ebu/mcma-libraries-go/client"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The url to the services endpoint of the MCMA Service Registry",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MCMA_SERVICES_URL", nil),
			},
			"auth": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "The type of authentication to use",
							Required:    true,
						},
						"data": {
							Type:        schema.TypeMap,
							Description: "Data used by this authentication type, e.g. keys, profile names, etc",
							Optional:    true,
						},
						"use_for_initialization": {
							Type:        schema.TypeBool,
							Description: "Indicates if this auth type should be used to initialize the provider with service data",
							Optional:    true,
						},
					},
				},
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
	if url == "" {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create MCMA ResourceManager",
				Detail:   "Url for MCMA service registry was not provided, and the MCMA_SERVICES_URL environment variable is not set.",
			},
		}
	}

	var initAuthType string
	authMap := make(map[string]mcma.Authenticator)

	for _, raw := range d.Get("auth").(*schema.Set).List() {
		auth := raw.(map[string]interface{})

		authType := auth["type"].(string)
		if ufi, found := auth["use_for_initialization"]; found && ufi.(bool) {
			initAuthType = authType
		}

		authenticator, d := GetAuthenticator(authType, auth)
		if d != nil {
			return nil, d
		}

		authMap[authType] = authenticator
	}

	var resourceManager mcma.ResourceManager
	if len(initAuthType) != 0 {
		resourceManager = mcma.NewResourceManager(url, initAuthType)
	} else {
		resourceManager = mcma.NewResourceManagerNoAuth(url)
	}

	for key, a := range authMap {
		resourceManager.AddAuth(key, a)
	}

	return &resourceManager, nil
}
