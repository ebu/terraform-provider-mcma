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

func TestAccMcmaJobProfile_basic(t *testing.T) {
	var jobProfile mcmamodel.JobProfile
	profileName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckMcmaJobProfileDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccountMcmaJobProfile(profileName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJobProfileExists("mcma_job_profile.job_profile_"+profileName, &jobProfile),
				),
			},
		},
	})
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

func testAccountMcmaJobProfile(profileName string) string {
	return fmt.Sprintf(`
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
  custom_properties = {
	customprop1 = "customprop1val"
	customprop2 = "customprop2val"
  }
}
`, profileName, profileName)
}

func testAccCheckJobProfileExists(resourceName string, jobProfile *mcmamodel.JobProfile) resource.TestCheckFunc {
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
		pImpl := p.(mcmamodel.JobProfile)
		*jobProfile = pImpl
		return nil
	}
}
