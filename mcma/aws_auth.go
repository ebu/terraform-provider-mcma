package mcma

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	mcmaclient "github.com/ebu/mcma-libraries-go/client"
)

func GetAWS4Authenticator(authData map[string]interface{}) (mcmaclient.Authenticator, diag.Diagnostics) {
	region, d := GetAuthDataString(authData, "region", false)
	if d != nil {
		return nil, d
	}

	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			return nil, diag.Errorf("region not specified and AWS_REGION environment variable not set")
		}
	}

	accessKey, d := GetAuthDataString(authData, "access_key", false)
	if d != nil {
		return nil, d
	}

	if len(accessKey) > 0 {
		secretKey, d := GetAuthDataString(authData, "secret_key", true)
		if d != nil {
			return nil, d
		}
		sessionToken, d := GetAuthDataString(authData, "session_token", false)
		if d != nil {
			return nil, d
		}
		return mcmaclient.NewAWS4AuthenticatorFromKeys(accessKey, secretKey, sessionToken, region), nil
	}

	profile, d := GetAuthDataString(authData, "profile", false)
	if d != nil {
		return nil, d
	}
	if len(profile) > 0 {
		return mcmaclient.NewAWS4AuthenticatorFromProfile(profile, region), nil
	}

	return mcmaclient.NewAWS4AuthenticatorFromEnvVars(), nil
}
