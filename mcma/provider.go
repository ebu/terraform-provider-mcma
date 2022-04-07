package mcma

import (
	"context"

	mcma "github.com/ebu/mcma-libraries-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"services_url": {
				Type:        schema.TypeString,
				Description: "The url to the services endpoint of the MCMA Service Registry",
				Required:    true,
			},
			"services_auth_type": {
				Type:        schema.TypeString,
				Description: "The auth type to use for the services endpoint of the MCMA Service Registry",
				Optional:    true,
			},
			"aws4_auth": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "The AWS region to use for authentication",
							Optional:    true,
						},
						"profile": {
							Type:        schema.TypeString,
							Description: "The AWS profile to use for authentication",
							Optional:    true,
						},
						"access_key": {
							Type:        schema.TypeString,
							Description: "The AWS access key to use for authentication. Requires that secret_key also be specified",
							Optional:    true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Description: "The AWS secret key to use for authentication. Requires that access_key also be specified",
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

func addAuthToMap(
	authMap map[string]mcma.Authenticator,
	resourceData *schema.ResourceData,
	authType string,
	authFactory func(map[string]interface{}) (mcma.Authenticator, diag.Diagnostics),
) diag.Diagnostics {
	blocks := resourceData.Get(authType + "_auth").(*schema.Set).List()
	switch len(blocks) {
	case 0:
		return nil
	case 1:
		authenticator, d := authFactory(blocks[0].(map[string]interface{}))
		if d != nil {
			return d
		}
		authMap[authType] = authenticator
		return nil
	default:
		return diag.Errorf("only 1 %s_auth block allowed", authType)
	}
}

func configure(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	servicesUrl := d.Get("services_url").(string)
	if servicesUrl == "" {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create MCMA ResourceManager",
				Detail:   "Url for MCMA service registry was not provided, and the MCMA_SERVICES_URL environment variable is not set.",
			},
		}
	}
	servicesAuthType := d.Get("services_auth_type").(string)

	authMap := make(map[string]mcma.Authenticator)
	addAuthToMap(authMap, d, "aws4", GetAWS4Authenticator)

	if len(authMap) == 1 && servicesAuthType == "" {
		for s := range authMap {
			servicesAuthType = s
		}
	}

	var resourceManager mcma.ResourceManager
	if len(servicesAuthType) != 0 {
		resourceManager = mcma.NewResourceManager(servicesUrl, servicesAuthType)
	} else {
		resourceManager = mcma.NewResourceManagerNoAuth(servicesUrl)
	}

	for key, a := range authMap {
		resourceManager.AddAuth(key, a)
	}

	return &resourceManager, nil
}
