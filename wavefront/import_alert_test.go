package wavefront_plugin

import (
	"testing"

	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func TestAccAlert_importThreshold(t *testing.T) {
	resourceName := "wavefront_alert.test_threshold_alert"
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontAlertImporter_threshold(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_threshold_alert", &record),
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

func testAccCheckWavefrontAlertImporter_threshold() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
  name = "Terraform Test Target Import"
  description = "Test target"
  method = "EMAIL"
  recipient = "test@example.com"
  email_subject = "This is a test"
  is_html_content = true
  template = "{}"
  triggers = [
    "ALERT_OPENED",
    "ALERT_RESOLVED"
  ]
}


resource "wavefront_alert" "test_threshold_alert" {
  name = "Terraform Test Alert Import"
  alert_type = "THRESHOLD"
  additional_information = "This is a Terraform Test Alert"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5

  threshold_conditions = {
    "severe" = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
    "warn" = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 60"
    "info" = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 50"
  }

  threshold_targets = {
	"severe" = "target:${wavefront_alert_target.test_target.id}"
  }
  
  tags = [
    "terraform"
  ]
}
`)
}
