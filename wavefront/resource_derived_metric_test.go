package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccWavefrontDerivedMetric_Basic(t *testing.T) {
	var record wavefront.DerivedMetric

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDerivedMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDerivedMetric_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDerivedMetricExists("wavefront_derived_metric.derived", &record),
					testAccCheckWavefrontDerivedMetricAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "name", "dummy derived metric"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "query", "aliasMetric(5, \"some.metric\")"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "additional_information", "this is a dummy derived metric"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "tags.#", "3"),
				),
			},
		},
	})
}

func TestAccWavefrontDerivedMetric_Updated(t *testing.T) {
	var record wavefront.DerivedMetric

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDerivedMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDerivedMetric_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDerivedMetricExists("wavefront_derived_metric.derived", &record),
					testAccCheckWavefrontDerivedMetricAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "minutes", "5"),
				),
			},
			{
				Config: testAccCheckWavefrontDerivedMetric_Updated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDerivedMetricExists("wavefront_derived_metric.derived", &record),
					testAccCheckWavefrontDerivedMetricAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "minutes", "10"),
				),
			},
		},
	})
}

func TestAccWavefrontDerivedMetric_Multiple(t *testing.T) {
	var record wavefront.DerivedMetric
	var record2 wavefront.DerivedMetric

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDerivedMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDerivedMetric_Multiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDerivedMetricExists("wavefront_derived_metric.derived", &record),
					testAccCheckWavefrontDerivedMetricAttributes(&record),
					testAccCheckWavefrontDerivedMetricExists("wavefront_derived_metric.derived", &record2),
					testAccCheckWavefrontDerivedMetricAttributes(&record2),

					// Check the first record
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "name", "dummy derived metric"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "query", "aliasMetric(5, \"some.metric\")"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "additional_information", "this is a dummy derived metric"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived", "tags.#", "3"),
					// Check the second record
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived2", "name", "dummy derived metric"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived2", "minutes", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived2", "query", "aliasMetric(5, \"some.metric\")"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived2", "additional_information", "this is a dummy derived metric"),
					resource.TestCheckResourceAttr(
						"wavefront_derived_metric.derived2", "tags.#", "3"),
				),
			},
		},
	})
}

func testAccCheckWavefrontDerivedMetricAttributes(dm *wavefront.DerivedMetric) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dm.Name != "dummy derived metric" {
			return fmt.Errorf("bad value: %s", dm.Name)
		}

		if !(dm.Minutes == 5 || dm.Minutes == 10) {
			return fmt.Errorf("bad value: %d", dm.Minutes)
		}

		if dm.Query != "aliasMetric(5, \"some.metric\")" {
			return fmt.Errorf("bad value: %s", dm.Query)
		}

		if len(dm.Tags.CustomerTags) != 3 {
			return fmt.Errorf("expected 3 tag values, got %d", len(dm.Tags.CustomerTags))
		}

		if dm.AdditionalInformation != "this is a dummy derived metric" {
			return fmt.Errorf("bad value: %s", dm.AdditionalInformation)
		}

		return nil
	}
}

func testAccCheckWavefrontDerivedMetricExists(n string, dm *wavefront.DerivedMetric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		alerts := testAccProvider.Meta().(*wavefrontClient).client.DerivedMetrics()
		tmp := wavefront.DerivedMetric{ID: &rs.Primary.ID}

		err := alerts.Get(&tmp)
		if err != nil {
			return fmt.Errorf("Error finding Wavefront DerivedMetric %s", err)
		}

		*dm = tmp

		return nil
	}
}

func testAccCheckWavefrontDerivedMetricDestroy(s *terraform.State) error {

	derivedMetrics := testAccProvider.Meta().(*wavefrontClient).client.DerivedMetrics()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_derived_metric" {
			continue
		}

		tmpDM := wavefront.DerivedMetric{ID: &rs.Primary.ID}

		err := derivedMetrics.Get(&tmpDM)
		if err == nil {
			return fmt.Errorf("DerivedMetric still exists")
		}
	}

	return nil
}

func testAccCheckWavefrontDerivedMetric_Basic() string {
	return fmt.Sprintf(`
resource "wavefront_derived_metric" "derived" {
  name                   = "dummy derived metric"
  minutes                = 5
  query                  = "aliasMetric(5, \"some.metric\")"
  additional_information = "this is a dummy derived metric"
  tags = [
    "somemetric",
    "thatistagged",
    "withmytags"
  ]
}
`)
}

func testAccCheckWavefrontDerivedMetric_Updated() string {
	return fmt.Sprintf(`
resource "wavefront_derived_metric" "derived" {
  name                   = "dummy derived metric"
  minutes                = 10
  query                  = "aliasMetric(5, \"some.metric\")"
  additional_information = "this is a dummy derived metric"
  tags = [
    "somemetric",
    "thatistagged",
    "withmytags"
  ]
}
`)
}

func testAccCheckWavefrontDerivedMetric_Multiple() string {
	return fmt.Sprintf(`
resource "wavefront_derived_metric" "derived" {
  name                   = "dummy derived metric"
  minutes                = 5
  query                  = "aliasMetric(5, \"some.metric\")"
  additional_information = "this is a dummy derived metric"
  tags = [
    "somemetric",
    "thatistagged",
    "withmytags"
  ]
}

resource "wavefront_derived_metric" "derived2" {
  name                   = "dummy derived metric"
  minutes                = 5
  query                  = "aliasMetric(5, \"some.metric\")"
  additional_information = "this is a dummy derived metric"
  tags = [
    "somemetric",
    "thatistagged",
    "withmytags"
  ]
}
`)
}
