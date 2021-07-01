package mcma

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider {
		"mcma": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MCMA_SERVICES_URL"); v == "" {
		t.Fatal("MCMA_SERVICES_URL must be set for acceptance tests")
	}
	if v := os.Getenv("MCMA_SERVICES_AUTH_TYPE"); v == "" {
		t.Fatal("MCMA_SERVICES_AUTH_TYPE must be set for acceptance tests")
	}
	if v := os.Getenv("MCMA_SERVICES_AUTH_CONTEXT"); v == "" {
		t.Fatal("MCMA_SERVICES_AUTH_CONTEXT must be set for acceptance tests")
	}
}
