package mcma

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
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

func getProviderConfig(servicesUrl string, servicesAuthType string, authBlocks []aws4AuthBlock) string {
	providerConfig := "provider \"mcma\" {\n"
	providerConfig += "  services_url = \"" + servicesUrl + "\"\n"
	if servicesAuthType != "" {
		providerConfig += "  services_auth_type = \"" + servicesAuthType + "\"\n"
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

func getAwsProfileProviderConfig(servicesUrl string, region string, profile string) string {
	authBlocks := make([]aws4AuthBlock, 1)
	authBlocks[0] = aws4AuthBlock{
		region:  region,
		profile: profile,
	}
	return getProviderConfig(servicesUrl, "", authBlocks)
}

func getAwsProfileProviderConfigFromEnvVars() string {
	return getAwsProfileProviderConfig(os.Getenv("MCMA_AWS_SERVICES_URL"), os.Getenv("MCMA_AWS_REGION"), os.Getenv("MCMA_AWS_PROFILE"))
}

func getKubernetesProviderConfigFromEnvVars() string {
	return getProviderConfig(os.Getenv("MCMA_KUBERNETES_SERVICES_URL"), "", nil)
}
