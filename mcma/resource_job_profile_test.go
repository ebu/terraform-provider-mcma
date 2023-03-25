package mcma

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	mcmaclient "github.com/ebu/mcma-libraries-go/client"
	mcmamodel "github.com/ebu/mcma-libraries-go/model"
)

func TestAccMcmaJobProfile_basic(t *testing.T) {
	profileName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	createTestCase := func(providerConfig string) resource.TestCase {
		return resource.TestCase{
			Providers: testAccProviders,
			CheckDestroy: resource.ComposeTestCheckFunc(
				testAccCheckMcmaJobProfileDestroy,
			),
			Steps: []resource.TestStep{
				{
					Config: testAccountMcmaJobProfile(profileName, providerConfig),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckJobProfileExists("mcma_job_profile.job_profile_" + profileName),
					),
				},
				{
					Config: testAccountMcmaJobProfileMultiple(profileName, providerConfig),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckJobProfileExists("mcma_job_profile.job_profile_"+profileName+"_1"),
						testAccCheckJobProfileExists("mcma_job_profile.job_profile_"+profileName+"_2"),
						testAccCheckJobProfileExists("mcma_job_profile.job_profile_"+profileName+"_3"),
					),
				},
			},
		}
	}
	//resource.Test(t, createTestCase(getKubernetesProviderConfigFromEnvVars()))
	resource.Test(t, createTestCase(getAwsProfileProviderConfigFromEnvVars()))
}

func testAccCheckMcmaJobProfileDestroy(s *terraform.State) error {
	resourceManager := testAccProvider.Meta().(*mcmaclient.ResourceManager)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mcma_job_profile" {
			continue
		}
		existing, err := resourceManager.Get(reflect.TypeOf(mcmamodel.JobProfile{}), rs.Primary.ID)
		if err != nil {
			return err
		}
		for i := 0; existing != nil && i < 30; i++ {
			time.Sleep(1 * time.Second)
			existing, err = resourceManager.Get(reflect.TypeOf(mcmamodel.JobProfile{}), rs.Primary.ID)
			if err != nil {
				return err
			}
		}
		if existing != nil {
			return fmt.Errorf("job profile (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccountMcmaJobProfile(profileName string, providerConfig string) string {
	return fmt.Sprintf(`
%s

resource "mcma_job_profile" "job_profile_%s" {
  name = "%s"
  input_parameter {
	name = "param1"
	type = "string"
  }
  input_parameter {
	name = "param2"
	type = "number"
	optional = true
  }
  output_parameter {
	name = "outparam1"
	type = "string"
  }
  output_parameter {
	name = "outparam2"
	type = "number"
  }
  custom = {
	customprop1 = "customprop1val"
	customprop2 = "customprop2val"
  }
}
`, providerConfig, profileName, profileName)
}

func testAccountMcmaJobProfileMultiple(profileName string, providerConfig string) string {
	return fmt.Sprintf(`
%s

resource "mcma_job_profile" "job_profile_%s_1" {
  name = "%s_1"
  input_parameter {
	name = "param1"
	type = "string"
  }
  input_parameter {
	name = "param2"
	type = "number"
	optional = true
  }
  output_parameter {
	name = "outparam1"
	type = "string"
  }
  output_parameter {
	name = "outparam2"
	type = "number"
  }
  custom = {
	customprop1 = "customprop1val"
	customprop2 = "customprop2val"
  }
}

resource "mcma_job_profile" "job_profile_%s_2" {
  name = "%s_2"
  input_parameter {
	name = "param3"
	type = "string"
  }
  input_parameter {
	name = "param4"
	type = "number"
	optional = true
  }
  output_parameter {
	name = "outparam3"
	type = "string"
  }
  output_parameter {
	name = "outparam4"
	type = "number"
  }
  custom = {
	customprop1 = "customprop3val"
	customprop2 = "customprop4val"
  }
}

resource "mcma_job_profile" "job_profile_%s_3" {
  name = "%s_3"
  input_parameter {
	name = "param5"
	type = "string"
  }
  input_parameter {
	name = "param6"
	type = "number"
	optional = true
  }
  output_parameter {
	name = "outparam5"
	type = "string"
  }
  output_parameter {
	name = "outparam6"
	type = "number"
  }
  custom = {
	customprop1 = "customprop5val"
	customprop2 = "customprop6val"
  }
}
`, providerConfig, profileName, profileName, profileName, profileName, profileName, profileName)
}

func testAccCheckJobProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("job profile ID not set")
		}
		resourceManager := testAccProvider.Meta().(*mcmaclient.ResourceManager)
		p, err := resourceManager.Get(reflect.TypeOf(mcmamodel.JobProfile{}), rs.Primary.ID)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("job profile with ID %s not found", rs.Primary.ID)
		}
		return nil
	}
}
