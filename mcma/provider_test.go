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

type authBlock struct {
	Type                 string
	Data                 map[string]string
	UseForInitialization bool
}

func getProviderConfig(url string, authBlocks []authBlock) string {
	providerConfig := "provider \"mcma\" {\n"
	providerConfig += "  url = \"" + url + "\"\n"
	if authBlocks != nil && len(authBlocks) > 0 {
		for _, authBlock := range authBlocks {
			providerConfig += "  auth {\n"
			providerConfig += "    type = \"" + authBlock.Type + "\"\n"
			if authBlock.Data != nil {
				providerConfig += "    data = {\n"
				for key, val := range authBlock.Data {
					providerConfig += "      " + key + " = \"" + val + "\"\n"
				}
				providerConfig += "    }\n"
			}
			if authBlock.UseForInitialization {
				providerConfig += "    use_for_initialization = true\n"
			}
			providerConfig += "  }\n"
		}
	}
	providerConfig += "}\n"
	return providerConfig
}

func getAwsProfileProviderConfig(url string, region string, profile string) string {
	awsData := make(map[string]string)
	awsData["region"] = region
	awsData["profile"] = profile
	authBlocks := make([]authBlock, 1)
	authBlocks[0] = authBlock{
		"aws4",
		awsData,
		true,
	}
	return getProviderConfig(url, authBlocks)
}

func getAwsProfileProviderConfigFromEnvVars() string {
	return getAwsProfileProviderConfig(os.Getenv("MCMA_AWS_SERVICES_URL"), os.Getenv("MCMA_AWS_REGION"), os.Getenv("MCMA_AWS_PROFILE"))
}

func getKubernetesProviderConfigFromEnvVars() string {
	return getProviderConfig(os.Getenv("MCMA_KUBERNETES_SERVICES_URL"), nil)
}
