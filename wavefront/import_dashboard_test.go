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
  display_section_table_of_contents = true
  display_query_parameters = true
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
        summarization = "MEAN"
        chart_setting {
          auto_column_tags = false
          column_tags = "deprecated"
          custom_tags = ["tag1", "tag2"]
          expected_data_spacing = 120
          fixed_legend_display_stats = ["stat1", "stat2"]
          fixed_legend_enabled = true
          fixed_legend_filter_field = "CURRENT"
          fixed_legend_filter_limit = 1
          fixed_legend_filter_sort = "TOP"
          fixed_legend_hide_label = false
          fixed_legend_position = "RIGHT"
          fixed_legend_use_raw_stats = true
          group_by_source = true
          num_tags = 2
          plain_markdown_content = "markdown content"
          tag_mode = "all"
          type = "markdown-widget"
          windowing = "full"
          window_size = 10
          max = 100
          min = 0
          y0_scale_si_by_1024 = true
          y0_unit_autoscaling = true
          y1max = 100
          y1_scale_si_by_1024 = true
          y1_unit_autoscaling = true
          y1_units = "units"
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
