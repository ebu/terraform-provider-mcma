package mcma

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"mcma": testAccProvider,
	}
}

type aws4AuthBlock struct {
	region    string
	profile   string
	accessKey string
	secretKey string
}

func getProviderConfig(serviceRegistryUrl string, serviceRegistryAuthType string, authBlocks []aws4AuthBlock) string {
	providerConfig := "provider \"mcma\" {\n"
	providerConfig += "  service_registry_url = \"" + serviceRegistryUrl + "\"\n"
	if serviceRegistryAuthType != "" {
		serviceRegistryAuthType += "  service_registry_auth_type = \"" + serviceRegistryAuthType + "\"\n"
	}
	if authBlocks != nil && len(authBlocks) > 0 {
		for _, authBlock := range authBlocks {
			providerConfig += "  aws4_auth {\n"
			if authBlock.region != "" {
				providerConfig += "    region = \"" + authBlock.region + "\"\n"
			}
			if authBlock.profile != "" {
				providerConfig += "    profile = \"" + authBlock.profile + "\"\n"
			}
			if authBlock.accessKey != "" {
				providerConfig += "    access_key = \"" + authBlock.accessKey + "\"\n"
			}
			if authBlock.secretKey != "" {
				providerConfig += "    secret_key = \"" + authBlock.secretKey + "\"\n"
			}
			providerConfig += "  }\n"
		}
	}
	providerConfig += "}\n"
	return providerConfig
}

func getAwsProfileProviderConfig(serviceRegistryUrl string, region string, profile string) string {
	authBlocks := make([]aws4AuthBlock, 1)
	authBlocks[0] = aws4AuthBlock{
		region:  region,
		profile: profile,
	}
	return getProviderConfig(serviceRegistryUrl, "", authBlocks)
}

func getAwsProfileProviderConfigFromEnvVars() string {
	return getAwsProfileProviderConfig(os.Getenv("MCMA_AWS_SERVICE_REGISTRY_URL"), os.Getenv("MCMA_AWS_REGION"), os.Getenv("MCMA_AWS_PROFILE"))
}

func getKubernetesProviderConfigFromEnvVars() string {
	return getProviderConfig(os.Getenv("MCMA_KUBERNETES_SERVICE_REGISTRY_URL"), "", nil)
}
