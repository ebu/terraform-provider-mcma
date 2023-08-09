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

type authBlock interface {
	GetText() string
}

type aws4AuthBlock struct {
	region    string
	profile   string
	accessKey string
	secretKey string
}

func (authBlock aws4AuthBlock) GetText() string {
	authBlockText := "  aws4_auth {\n"
	if authBlock.region != "" {
		authBlockText += "    region = \"" + authBlock.region + "\"\n"
	}
	if authBlock.profile != "" {
		authBlockText += "    profile = \"" + authBlock.profile + "\"\n"
	}
	if authBlock.accessKey != "" {
		authBlockText += "    access_key = \"" + authBlock.accessKey + "\"\n"
	}
	if authBlock.secretKey != "" {
		authBlockText += "    secret_key = \"" + authBlock.secretKey + "\"\n"
	}
	authBlockText += "  }\n"
	return authBlockText
}

type mcmaApiKeyAuthBlock struct {
	apiKey string
}

func (authBlock mcmaApiKeyAuthBlock) GetText() string {
	authBlockText := "  mcma_api_key_auth {\n"
	if authBlock.apiKey != "" {
		authBlockText += "    api_key = \"" + authBlock.apiKey + "\"\n"
	}
	authBlockText += "  }\n"
	return authBlockText
}

func getProviderConfig(serviceRegistryUrl string, serviceRegistryAuthType string, authBlocks []authBlock) string {
	providerConfig := "provider \"mcma\" {\n"
	providerConfig += "  service_registry_url = \"" + serviceRegistryUrl + "\"\n"
	if serviceRegistryAuthType != "" {
		serviceRegistryAuthType += "  service_registry_auth_type = \"" + serviceRegistryAuthType + "\"\n"
	}
	if authBlocks != nil && len(authBlocks) > 0 {
		for _, authBlock := range authBlocks {
			providerConfig += authBlock.GetText()
		}
	}
	providerConfig += "}\n"
	return providerConfig
}

func getAwsProfileProviderConfig(serviceRegistryUrl string, region string, profile string) string {
	authBlocks := make([]authBlock, 1)
	authBlocks[0] = aws4AuthBlock{
		region:  region,
		profile: profile,
	}
	return getProviderConfig(serviceRegistryUrl, "", authBlocks)
}

func getAwsProfileProviderConfigFromEnvVars() string {
	return getAwsProfileProviderConfig(os.Getenv("MCMA_AWS_SERVICE_REGISTRY_URL"), os.Getenv("MCMA_AWS_REGION"), os.Getenv("MCMA_AWS_PROFILE"))
}

func getMcmaApiKeyProviderConfig(serviceRegistryUrl, apiKey string) string {
	authBlocks := make([]authBlock, 1)
	authBlocks[0] = mcmaApiKeyAuthBlock{
		apiKey: apiKey,
	}
	return getProviderConfig(serviceRegistryUrl, "", authBlocks)
}

func getMcmaApiKeyProviderConfigFromEnvVars() string {
	return getMcmaApiKeyProviderConfig(os.Getenv("MCMA_AWS_SERVICE_REGISTRY_URL"), os.Getenv("MCMA_API_KEY"))
}
