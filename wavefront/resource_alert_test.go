package wavefront_plugin

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spaceapegames/go-wavefront"
)

func TestAccWavefrontAlert_Basic(t *testing.T) {
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckWavefrontAlert_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert", &record),
					testAccCheckWavefrontAlertAttributes(&record),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "name", "Terraform Test Alert"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "target", "test@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "condition", "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "display_expression", "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "resolve_after_minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "severity", "WARN"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "tags.0", "terraform"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "tags.1", "test"),
				),
			},
		},
	})
}

func TestAccWavefrontAlert_Updated(t *testing.T) {
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontAlert_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert", &record),
					testAccCheckWavefrontAlertAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "target", "test@example.com"),
				),
			},
			{
				Config: testAccCheckWavefrontAlert_new_value(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert", &record),
					testAccCheckWavefrontAlertAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "target", "terraform@example.com"),
				),
			},
		},
	})
}

// Fails due to Wavefront known issue. Uncomment when that is fixed.
//func TestAccWavefrontAlert_Multiple(t *testing.T) {
//	var record wavefront.Alert
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckWavefrontAlertDestroy,
//		Steps: []resource.TestStep{
//			resource.TestStep{
//				Config: testAccCheckWavefrontAlert_multiple(),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert1", &record),
//					testAccCheckWavefrontAlertAttributes(&record),
//					resource.TestCheckResourceAttr(
//						"wavefront_alert.test_alert1", "name", "Terraform Test Alert1"),
//					resource.TestCheckResourceAttr(
//						"wavefront_alert.test_alert2", "name", "Terraform Test Alert2"),
//					resource.TestCheckResourceAttr(
//						"wavefront_alert.test_alert3", "name", "Terraform Test Alert3"),
//				),
//			},
//		},
//	})
//}

func testAccCheckWavefrontAlertDestroy(s *terraform.State) error {

	alerts := testAccProvider.Meta().(*wavefrontClient).client.Alerts()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_alert" {
			continue
		}

		results, err := alerts.Find(
			[]*wavefront.SearchCondition{
				&wavefront.SearchCondition{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("Error finding Wavefront Alert. %s", err)
		}
		if len(results) > 0 {
			return fmt.Errorf("Alert still exists")
		}
	}

	return nil
}

func testAccCheckWavefrontAlertAttributes(alert *wavefront.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Target != "test@example.com" {
			return fmt.Errorf("Bad value: %s", alert.Target)
		}

		return nil
	}
}

func testAccCheckWavefrontAlertAttributesUpdated(alert *wavefront.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Target != "terraform@example.com" {
			return fmt.Errorf("Bad value: %s", alert.Target)
		}

		return nil
	}
}

func testAccCheckWavefrontAlertExists(n string, alert *wavefront.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		alerts := testAccProvider.Meta().(*wavefrontClient).client.Alerts()

		results, err := alerts.Find(
			[]*wavefront.SearchCondition{
				&wavefront.SearchCondition{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("Error finding Wavefront Alert %s", err)
		}
		// resource has been deleted out of band. So unset ID
		if len(results) != 1 {
			return fmt.Errorf("No Alerts Found")
		}
		if *results[0].ID != rs.Primary.ID {
			return fmt.Errorf("Alert not found")
		}

		*alert = *results[0]

		return nil
	}
}

func testAccCheckWavefrontAlert_basic() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "test_alert" {
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

func testAccCheckWavefrontAlert_new_value() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "test_alert" {
  name = "Terraform Test Alert"
  target = "terraform@example.com"
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

func testAccCheckWavefrontAlert_multiple() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "test_alert1" {
  name = "Terraform Test Alert 1"
  target = "terraform@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform1",
  ]
}
resource "wavefront_alert" "test_alert2" {
  name = "Terraform Test Alert 2"
  target = "terraform@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform2",
    "test"
  ]
}
resource "wavefront_alert" "test_alert3" {
  name = "Terraform Test Alert 3"
  target = "terraform@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform",
  ]
}
`)
}
