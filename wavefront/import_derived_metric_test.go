package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDerivedMetric_importBasic(t *testing.T) {
	resourceName := "wavefront_derived_metric.derived"
	var record wavefront.DerivedMetric

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDerivedMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDerivedMetricImporter_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDerivedMetricExists("wavefront_derived_metric.derived", &record),
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

func testAccCheckWavefrontDerivedMetricImporter_Basic() string {
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
