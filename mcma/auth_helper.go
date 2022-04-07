package mcma

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"reflect"
)

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
