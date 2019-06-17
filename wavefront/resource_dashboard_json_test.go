package wavefront_plugin

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spaceapegames/go-wavefront"
	"testing"
)

func TestAccWavefrontDashboardJson_Basic(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardJsonDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboardJson_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJsonExists("wavefront_dashboard_json.test_dashboard_json", &record),
					testAccCheckWavefrontDashboardJsonAttributes(&record),

					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_json", "id", "tftestimport"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboardJson_Updated(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardJsonDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboardJson_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJsonExists("wavefront_dashboard_json.test_dashboard_json", &record),
					testAccCheckWavefrontDashboardJsonAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_json", "id", "tftestimport"),
				),
			},
			{
				Config: testAccCheckWavefrontDashboardJson_new_value(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJsonExists("wavefront_dashboard_json.test_dashboard_json", &record),
					testAccCheckWavefrontDashboardJsonAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_json", "id", "tftestimport"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboardJson_Multiple(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardJsonDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboardJson_multiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJsonExists("wavefront_dashboard_json.test_dashboard_1", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_1", "id", "test_dashboard_1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_2", "id", "test_dashboard_2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_3", "id", "test_dashboard_3"),
				),
			},
		},
	})
}

func testAccCheckWavefrontDashboardJsonDestroy(s *terraform.State) error {

	dashboards := testAccProvider.Meta().(*wavefrontClient).client.Dashboards()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_dashboard" {
			continue
		}

		results, err := dashboards.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("Error finding Wavefront Dashboard. %s", err)
		}
		if len(results) > 0 {
			return fmt.Errorf("Dashboard still exists")
		}
	}

	return nil
}

func testAccCheckWavefrontDashboardJsonAttributes(dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dashboard.Name != "Terraform Test Dashboard Json" {
			return fmt.Errorf("Bad value: %s", dashboard.Name)
		}

		return nil
	}
}

func testAccCheckWavefrontDashboardJsonAttributesUpdated(dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dashboard.Name != "Terraform Test Dashboard Json Updated" {
			return fmt.Errorf("Bad value: %s", dashboard.Name)
		}

		return nil
	}
}

func testAccCheckWavefrontDashboardJsonExists(n string, dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		dash := wavefront.Dashboard{
			ID: rs.Primary.ID,
		}

		dashboards := testAccProvider.Meta().(*wavefrontClient).client.Dashboards()
		err := dashboards.Get(&dash)
		if err != nil {
			return fmt.Errorf("Did not find Dashboard with id %s, %s", rs.Primary.ID, err)
		}
		*dashboard = dash
		return nil
	}
}

func testAccCheckWavefrontDashboardJson_basic() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard_json" "test_dashboard_json" {
	dashboard_json = <<EOF
{
  "name": "Terraform Test Dashboard Json",
  "description": "a",
  "eventFilterType": "BYCHART",
  "eventQuery": "",
  "defaultTimeWindow": "",
  "url": "tftestimport",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "sections": [
    {
      "name": "section 1",
      "rows": [
        {
          "charts": [
            {
              "name": "chart 1",
              "sources": [
                {
                  "name": "source 1",
                  "query": "ts()",
                  "scatterPlotSource": "Y",
                  "querybuilderEnabled": false,
                  "sourceDescription": ""
                }
              ],
              "units": "someunit",
              "base": 0,
              "noDefaultEvents": false,
              "interpolatePoints": false,
              "includeObsoleteMetrics": false,
              "description": "This is chart 1, showing something",
              "chartSettings": {
                "type": "markdown-widget",
                "max": 100,
                "expectedDataSpacing": 120,
                "windowing": "full",
                "windowSize": 10,
                "autoColumnTags": false,
                "columnTags": "deprecated",
                "tagMode": "all",
                "numTags": 2,
                "customTags": [
                  "tag1",
                  "tag2"
                ],
                "groupBySource": true,
                "y1Max": 100,
                "y1Units": "units",
                "y0ScaleSIBy1024": true,
                "y1ScaleSIBy1024": true,
                "y0UnitAutoscaling": true,
                "y1UnitAutoscaling": true,
                "fixedLegendEnabled": true,
                "fixedLegendUseRawStats": true,
                "fixedLegendPosition": "RIGHT",
                "fixedLegendDisplayStats": [
                  "stat1",
                  "stat2"
                ],
                "fixedLegendFilterSort": "TOP",
                "fixedLegendFilterLimit": 1,
                "fixedLegendFilterField": "CURRENT",
                "plainMarkdownContent": "markdown content"
              },
              "summarization": "MEAN"
            }
          ],
          "heightFactor": 50
        }
      ]
    }
  ],
  "parameterDetails": {
    "param": {
      "hideFromView": false,
      "description": null,
      "allowAll": null,
      "tagKey": null,
      "queryValue": null,
      "dynamicFieldType": null,
      "reverseDynSort": null,
      "parameterType": "SIMPLE",
      "label": "test",
      "defaultValue": "Label",
      "valuesToReadableStrings": {
        "Label": "test"
      },
      "selectedLabel": "Label",
      "value": "test"
    }
  },
  "tags" :{
    "customerTags":  ["terraform"]
  }
}
EOF
}
`)
}

func testAccCheckWavefrontDashboardJson_new_value() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard_json" "test_dashboard_json" {
	dashboard_json = <<EOF
{
  "name": "Terraform Test Dashboard Json Updated",
  "description": "a",
  "eventFilterType": "BYCHART",
  "eventQuery": "",
  "defaultTimeWindow": "",
  "url": "tftestimport",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "sections": [
    {
      "name": "section 1",
      "rows": [
        {
          "charts": [
            {
              "name": "chart 1",
              "sources": [
                {
                  "name": "source 1",
                  "query": "ts()",
                  "scatterPlotSource": "Y",
                  "querybuilderEnabled": false,
                  "sourceDescription": ""
                }
              ],
              "units": "someunit",
              "base": 0,
              "noDefaultEvents": false,
              "interpolatePoints": false,
              "includeObsoleteMetrics": false,
              "description": "This is chart 1, showing something",
              "chartSettings": {
                "type": "markdown-widget",
                "max": 100,
                "expectedDataSpacing": 120,
                "windowing": "full",
                "windowSize": 10,
                "autoColumnTags": false,
                "columnTags": "deprecated",
                "tagMode": "all",
                "numTags": 2,
                "customTags": [
                  "tag1",
                  "tag2"
                ],
                "groupBySource": true,
                "y1Max": 100,
                "y1Units": "units",
                "y0ScaleSIBy1024": true,
                "y1ScaleSIBy1024": true,
                "y0UnitAutoscaling": true,
                "y1UnitAutoscaling": true,
                "fixedLegendEnabled": true,
                "fixedLegendUseRawStats": true,
                "fixedLegendPosition": "RIGHT",
                "fixedLegendDisplayStats": [
                  "stat1",
                  "stat2"
                ],
                "fixedLegendFilterSort": "TOP",
                "fixedLegendFilterLimit": 1,
                "fixedLegendFilterField": "CURRENT",
                "plainMarkdownContent": "markdown content"
              },
              "summarization": "MEAN"
            }
          ],
          "heightFactor": 50
        }
      ]
    }
  ],
  "parameterDetails": {
    "param": {
      "hideFromView": false,
      "parameterType": "SIMPLE",
      "label": "test",
      "defaultValue": "Label",
      "valuesToReadableStrings": {
        "Label": "test"
      },
      "selectedLabel": "Label",
      "value": "test"
    }
  },
  "tags" :{
    "customerTags":  ["terraform"]
  }
}
EOF
}
`)
}

func testAccCheckWavefrontDashboardJson_multiple() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard_json" "test_dashboard_1" {
  dashboard_json = <<EOF
{
  "name": "test_dashboard_1",
  "eventFilterType": "BYCHART",
  "url": "test_dashboard_1",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "sections": [
    {
      "name": "New Section",
      "rows": []
    }
  ],
  "parameterDetails": {}
}
EOF
}
resource "wavefront_dashboard_json" "test_dashboard_2" {
  dashboard_json = <<EOF
{
  "name": "test_dashboard_2",
  "eventFilterType": "BYCHART",
  "url": "test_dashboard_2",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "sections": [
    {
      "name": "New Section",
      "rows": []
    }
  ],
  "parameterDetails": {}
}
EOF
}
resource "wavefront_dashboard_json" "test_dashboard_3" {
  dashboard_json = <<EOF
{
  "name": "test_dashboard_3",
  "eventFilterType": "BYCHART",
  "url": "test_dashboard_3",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "sections": [
    {
      "name": "New Section",
      "rows": []
    }
  ],
  "parameterDetails": {}
}
EOF
}
`)
}
