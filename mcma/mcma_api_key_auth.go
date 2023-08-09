package mcma

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	mcmaclient "github.com/ebu/mcma-libraries-go/client"
)

func GetMcmaApiKeyAuthenticator(authData map[string]interface{}) (mcmaclient.Authenticator, diag.Diagnostics) {
	apiKey, d := GetAuthDataString(authData, "api_key", true)
	if d != nil {
		return nil, d
	}

	return mcmaclient.NewMcmaApiKeyAuthenticator(apiKey), nil
}
