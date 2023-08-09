package mcma

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	mcmaclient "github.com/ebu/mcma-libraries-go/client"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"service_registry_url": {
				Type:        schema.TypeString,
				Description: "The url to the services endpoint of the MCMA Service Registry",
				Required:    true,
			},
			"service_registry_auth_type": {
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
			"mcma_api_key_auth": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:        schema.TypeString,
							Description: "The MCMA API key (header = 'x-mcma-api-key') to use for authentication",
							Required:    true,
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
	authMap map[string]mcmaclient.Authenticator,
	resourceData *schema.ResourceData,
	authType string,
	authKey string,
	authFactory func(map[string]interface{}) (mcmaclient.Authenticator, diag.Diagnostics),
) diag.Diagnostics {
	blocks := resourceData.Get(authKey + "_auth").(*schema.Set).List()
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
	serviceRegistryUrl := d.Get("service_registry_url").(string)
	if serviceRegistryUrl == "" {
		return nil, nil
	}
	serviceRegistryAuthType := d.Get("service_registry_auth_type").(string)

	authMap := make(map[string]mcmaclient.Authenticator)
	addAuthToMap(authMap, d, "AWS4", "aws4", GetAWS4Authenticator)
	addAuthToMap(authMap, d, "McmaApiKey", "mcma_api_key", GetMcmaApiKeyAuthenticator)

	if len(authMap) == 1 && serviceRegistryAuthType == "" {
		for s := range authMap {
			serviceRegistryAuthType = s
		}
	}

	var resourceManager mcmaclient.ResourceManager
	if len(serviceRegistryAuthType) != 0 {
		resourceManager = mcmaclient.NewResourceManager(serviceRegistryUrl, serviceRegistryAuthType)
	} else {
		resourceManager = mcmaclient.NewResourceManagerNoAuth(serviceRegistryUrl)
	}

	for key, a := range authMap {
		resourceManager.AddAuth(key, a)
	}

	return &resourceManager, nil
}

func getResourceManager(m interface{}) (*mcmaclient.ResourceManager, diag.Diagnostics) {
	if m == nil {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "ResourceManager not initialized.",
				Detail:   "Url for MCMA service registry was not provided, and the MCMA_SERVICE_REGISTRY_URL environment variable is not set.",
			},
		}
	}
	return m.(*mcmaclient.ResourceManager), nil
}
