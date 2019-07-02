package wavefront_plugin

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
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
			{
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
						"wavefront_alert.test_alert", "additional_information", "This is a Terraform Test Alert"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "display_expression", "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "resolve_after_minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "severity", "WARN"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert", "tags.#", "5"),
				),
			},
		},
	})
}

func TestAccWavefrontAlert_RequiredAttributes(t *testing.T) {
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontAlert_requiredAttributes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert_required", &record),
					testAccCheckWavefrontAlertAttributes(&record),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "name", "Terraform Test Alert Required Attributes Only"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "target", "test@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "condition", "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "additional_information", "This is a Terraform Test Alert Required"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "severity", "WARN"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "tags.#", "2"),
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

func TestAccWavefrontAlert_RemoveOptionalAttribute(t *testing.T) {
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontAlert_RemoveAttributes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert_required", &record),
					testAccCheckWavefrontAlertAttributesRemoved(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "target", "test@example.com"),
				),
			},
			{
				Config: testAccCheckWavefrontAlert_UpdatedRemoveAttributes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert_required", &record),
					testAccCheckWavefrontAlertAttributesRemovedUpdated(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert_required", "target", "terraform@example.com"),
				),
			},
		},
	})
}

func TestAccWavefrontAlert_Multiple(t *testing.T) {
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontAlert_multiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_alert1", &record),
					testAccCheckWavefrontAlertAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert1", "name", "Terraform Test Alert 1"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert2", "name", "Terraform Test Alert 2"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_alert3", "name", "Terraform Test Alert 3"),
				),
			},
		},
	})
}

func TestAccWavefrontAlert_Threshold(t *testing.T) {
	var record wavefront.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontAlert_threshold(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontAlertExists("wavefront_alert.test_threshold_alert", &record),
					testAccCheckWavefrontThresholdAlertAttributes(&record),

					//Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_threshold_alert", "threshold_conditions.%", "3"),
					resource.TestCheckResourceAttr(
						"wavefront_alert.test_threshold_alert", "threshold_targets.%", "1"),
				),
			},
		},
	})
}

func TestResourceAlert_validateAlertConditions(t *testing.T) {

	cases := []struct {
		name         string
		conf         *schema.ResourceData
		errorMessage string
	}{
		{
			"invalid alert type",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "WRONG")
				return d
			}(),
			"alert_type must be CLASSIC or THRESHOLD",
		},
		{
			"classic alert missing condition",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "CLASSIC")
				d.Set("severity", "severe")
				return d
			}(),
			"condition must be supplied for classic alerts",
		},
		{
			"classic alert missing severity",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "CLASSIC")
				d.Set("condition", "ts()")
				return d
			}(),
			"severity must be supplied for classic alerts",
		},
		{
			"classic alert",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "CLASSIC")
				d.Set("condition", "ts()")
				d.Set("severity", "severe")
				return d
			}(),
			"",
		},
		{
			"threshold alert missing conditions",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "THRESHOLD")
				return d
			}(),
			"threshold_conditions must be supplied for threshold alerts",
		},
		{
			"threshold alert",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "THRESHOLD")
				d.Set("threshold_conditions", map[string]interface{}{"severe": "ts()"})
				return d
			}(),
			"",
		},
		{
			"threshold alert invalid condition",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "THRESHOLD")
				d.Set("threshold_conditions", map[string]interface{}{"banana": "ts()"})
				return d
			}(),
			"invalid severity: banana",
		},
		{
			"threshold alert invalid target",
			func() *schema.ResourceData {
				d := resourceAlert().TestResourceData()
				d.Set("alert_type", "THRESHOLD")
				d.Set("threshold_conditions", map[string]interface{}{"severe": "ts()"})
				d.Set("threshold_targets", map[string]interface{}{"banana": "ts()"})

				return d
			}(),
			"invalid severity: banana",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateAlertConditions(&wavefront.Alert{}, c.conf)

			m := ""
			if err == nil {
				m = ""
			} else {
				m = err.Error()
			}

			if m != c.errorMessage {
				t.Errorf("expected error '%s', got '%s'", c.errorMessage, err.Error())
			}
		})
	}
}

func testAccCheckWavefrontAlertDestroy(s *terraform.State) error {

	alerts := testAccProvider.Meta().(*wavefrontClient).client.Alerts()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_alert" {
			continue
		}

		tmpAlert := wavefront.Alert{ID: &rs.Primary.ID}

		err := alerts.Get(&tmpAlert)
		if err == nil {
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

func testAccCheckWavefrontThresholdAlertAttributes(alert *wavefront.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if val, ok := alert.Conditions["severe"]; ok {
			if val != "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80" {
				return fmt.Errorf("bad value: %s", alert.Conditions["severe"])
			}
		} else {
			return fmt.Errorf("target not set")
		}

		return nil
	}
}

func testAccCheckWavefrontAlertAttributesUpdated(alert *wavefront.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Target != "terraform@example.com" {
			return fmt.Errorf("bad value: %s", alert.Target)
		}

		return nil
	}
}

func testAccCheckWavefrontAlertAttributesRemoved(alert *wavefront.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.ResolveAfterMinutes != 5 {
			return fmt.Errorf("unexpected value for ResolveAfterMinutes %v, expected 5", alert.ResolveAfterMinutes)
		}

		return nil
	}
}

func testAccCheckWavefrontAlertAttributesRemovedUpdated(alert *wavefront.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.ResolveAfterMinutes != 0 {
			return fmt.Errorf("unexpected value for ResolveAfterMinutes %v, expected 0", alert.ResolveAfterMinutes)
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
		tmpAlert := wavefront.Alert{ID: &rs.Primary.ID}

		err := alerts.Get(&tmpAlert)
		if err != nil {
			return fmt.Errorf("Error finding Wavefront Alert %s", err)
		}

		*alert = tmpAlert

		return nil
	}
}

func testAccCheckWavefrontAlert_basic() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "test_alert" {
  name = "Terraform Test Alert"
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
	"b",
    "terraform",
    "c",
    "test",
    "a"
  ]
}
`)
}

func testAccCheckWavefrontAlert_RemoveAttributes() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "test_alert_required" {
  name = "Terraform Test Alert Required Attributes Only"
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert Required"
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

func testAccCheckWavefrontAlert_UpdatedRemoveAttributes() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "test_alert_required" {
  name = "Terraform Test Alert Required Attributes Only"
  target = "terraform@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert Required"
  minutes = 5
  severity = "WARN"
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontAlert_requiredAttributes() string {
	return fmt.Sprintf(`
resource "wavefront_alert" "test_alert_required" {
  name = "Terraform Test Alert Required Attributes Only"
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert Required"
  minutes = 5
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
  additional_information = "This is a Terraform Test Alert"
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
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert"
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
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert"
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
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert"
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

func testAccCheckWavefrontAlert_threshold() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
  name = "Terraform Test Target"
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
  name = "Terraform Test Alert"
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
