package mcma

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"reflect"
	"strings"

	mcma "github.com/ebu/mcma-libraries-go/client"
)

func GetAuthenticator(authType string, auth map[string]interface{}) (mcma.Authenticator, diag.Diagnostics) {
	var authData map[string]interface{}
	if ad, found := auth["data"]; found {
		authData = *ad.(*map[string]interface{})
	}

	switch strings.ToLower(authType) {
	case "aws4":
		return GetAWS4Authenticator(authData)
	default:
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unsupported auth type",
				Detail:   fmt.Sprintf("The provider contains an 'auth' block that specifies unsupported auth type '%s'", authType),
			},
		}
	}
}

func GetAuthDataString(authData map[string]interface{}, key string, required bool) (string, diag.Diagnostics) {
	var value string
	if v, valFound := authData[key]; !valFound {
		if required {
			return "", diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("'%s' not specified in auth data", key),
					Detail:   fmt.Sprintf("A property with name '%s' must be specified for this authentication type", key),
				},
			}
		} else {
			return "", nil
		}
	} else {
		switch v.(type) {
		case string:
			value = v.(string)
			break
		default:
			return "", diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("'%s' be a string", key),
					Detail:   fmt.Sprintf("Expected a string value for '%s' in auth data but got a value of type", reflect.TypeOf(v).String()),
				},
			}
		}

		return value, nil
	}
}
