package wavefront_plugin

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/spaceapegames/go-wavefront"
)

func TestAccAlert_importBasic(t *testing.T) {
	resourceName := "wavefront_alert.foobar"
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontAlertImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.foobar", &record),
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

func testAccCheckWavefrontAlertImporter_basic() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "foobar" {
  name = "Terraform Test Alert"
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform",
    "test"
  ]
}
`)
}
