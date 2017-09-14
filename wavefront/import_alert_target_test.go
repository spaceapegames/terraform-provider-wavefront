package wavefront_plugin

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/spaceapegames/go-wavefront"
)

func TestAccTarget_importBasic(t *testing.T) {
	resourceName := "wavefront_alert_target.foobar"
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontTargetImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.foobar", &record),
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

func testAccCheckWavefrontTargetImporter_basic() string {
	return fmt.Sprintf(`
	resource "wavefront_alert_target" "foobar" {
	  name = "Terraform Test Target"
		description = "Test target"
		method = "WEBHOOK"
		recipient = "https://hooks.slack.com/services/test/me"
		content_type = "application/json"
		custom_headers = {
			"Testing" = "true"
		}
		template = "{}"
		triggers = [
			"ALERT_OPENED",
			"ALERT_RESOLVED"
		]
	}
	`)
}
