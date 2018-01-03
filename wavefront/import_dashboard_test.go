package wavefront_plugin

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/spaceapegames/go-wavefront"
)

func TestAccDashboard_importBasic(t *testing.T) {
	resourceName := "wavefront_dashboard.foobar"
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboardImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.foobar", &record),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckWavefrontDashboardImporter_basic() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "foobar" {
  name        = "Terraform Test Dashboard"
  description = "a"
  url         = "tftestimport"

  section {
    name = "section 1"

    row {
      chart {
        name        = "chart 1"
        description = "This is chart 1, showing something"
        units       = "someunit"

        source {
          name  = "source 1"
          query = "ts()"
        }
      }
    }
  }

  parameter_details {
    name = "param"
    label = "test"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "SIMPLE"
    values_to_readable_strings = {
      Label = "test"
    }
  }

  tags = [
    "terraform",
    "flamingo",
  ]
}
`)
}
