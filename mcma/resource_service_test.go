package mcma

import (
	"fmt"
	mcmaclient "github.com/ebu/mcma-libraries-go/client"
	mcmamodel "github.com/ebu/mcma-libraries-go/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"reflect"
	"testing"
	"time"
)

func TestAccMcmaService_basic(t *testing.T) {
	var service mcmamodel.Service
	profileName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckMcmaServiceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccountMcmaService(profileName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceExists("mcma_service.service_"+profileName, &service),
				),
			},
		},
	})
}

func testAccCheckMcmaServiceDestroy(s *terraform.State) error {
	resourceManager := testAccProvider.Meta().(*mcmaclient.ResourceManager)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mcma_service" {
			continue
		}
		existing, err := resourceManager.Get(reflect.TypeOf(mcmamodel.Service{}), rs.Primary.ID)
		if err != nil {
			return err
		}
		for i := 0; existing != nil && i < 30; i++ {
			time.Sleep(1 * time.Second)
			existing, err = resourceManager.Get(reflect.TypeOf(mcmamodel.Service{}), rs.Primary.ID)
			if err != nil {
				return err
			}
		}
		if existing != nil {
			return fmt.Errorf("service (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccountMcmaService(serviceName string) string {
	return fmt.Sprintf(`
resource "mcma_service" "service_%s" {
  name = "%s"
  auth_type = "AWS4"
  auth_context = "{ \"accessKey\": \"testkey\", \"secretKey\": \"testkey\", \"region\": \"us-east-1\" }"
  job_type = "AmeJob"
  resource {
	resource_type = "JobAssignment"
	http_endpoint = "https://some.endpoint.com/api/job-assignments"
  }
  resource {
	resource_type = "JobAssignment"
	http_endpoint = "https://some.endpoint.com/api/job-assignments"
	auth_type = "JWT"
    auth_context = "{ \"clientId\": \"some-client\" }"
  }
  input_location {
    url = "https://bucket.s3.amazonaws.com/prefix/input"
  }
  output_location {
    url = "https://bucket.s3.amazonaws.com/prefix/output"
  }
  job_profile_ids = [
     "https://service.registry.com/api/job-profiles/12345",
     "https://service.registry.com/api/job-profiles/67890"
  ]
}
`, serviceName, serviceName)
}

func testAccCheckServiceExists(resourceName string, service *mcmamodel.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("service ID not set")
		}
		resourceManager := testAccProvider.Meta().(*mcmaclient.ResourceManager)
		p, err := resourceManager.Get(reflect.TypeOf(mcmamodel.Service{}), rs.Primary.ID)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("service with ID %s not found", rs.Primary.ID)
		}
		pImpl := p.(mcmamodel.Service)
		*service = pImpl
		return nil
	}
}
