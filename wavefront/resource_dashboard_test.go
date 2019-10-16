package wavefront_plugin

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spaceapegames/go-wavefront"
)

func TestBuildTerraformParameterDetail(t *testing.T) {
	parameterDetail := wavefront.ParameterDetail{
		Label:         "test",
		DefaultValue:  "Label",
		HideFromView:  false,
		ParameterType: "SIMPLE",
		ValuesToReadableStrings: map[string]string{
			"Label": "test",
		},
	}

	result := buildTerraformParameterDetail(parameterDetail, "test")
	if result["label"] != "test" {
		t.Errorf("expected test got %s", result["label"])
	}
	if result["name"] != "test" {
		t.Errorf("expected test got %s", result["name"])
	}
}

func TestBuildTerraformSection(t *testing.T) {
	section := wavefront.Section{
		Name: "test",
		Rows: []wavefront.Row{},
	}

	result := buildTerraformSection(section)
	if result["name"] != "test" {
		t.Errorf("expected test, got %s", result["name"])
	}
	if len(result["row"].([]map[string]interface{})) != 0 {
		t.Errorf("Expected empty array, got Array of lenth %d", len(result["rows"].([]map[string]interface{})))
	}

	sectionWithRows := wavefront.Section{
		Rows: []wavefront.Row{
			{
				Charts: []wavefront.Chart{},
			},
			{
				Charts: []wavefront.Chart{},
			},
		},
	}

	resultWithRows := buildTerraformSection(sectionWithRows)
	if len(resultWithRows["row"].([]map[string]interface{})) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(resultWithRows["row"].([]map[string]interface{})))
	}
}

func TestBuildTerraformRow(t *testing.T) {
	row := wavefront.Row{
		Charts: []wavefront.Chart{},
	}

	result := buildTerraformRow(row)

	if len(result["chart"].([]map[string]interface{})) != 0 {
		t.Errorf("Expected empty array, got Array of lenth %d", len(result["chart"].([]map[string]interface{})))
	}

	rowWithCharts := wavefront.Row{
		Charts: []wavefront.Chart{
			{
				Name:    "chart 1",
				Sources: []wavefront.Source{},
			},
			{
				Name:    "chart 2",
				Sources: []wavefront.Source{},
			},
		},
	}

	resultWithCharts := buildTerraformRow(rowWithCharts)
	if len(resultWithCharts["chart"].([]map[string]interface{})) != 2 {
		t.Errorf("Expected array of length 2, got Array of lenth %d", len(result["chart"].([]map[string]interface{})))
	}
}

func TestBuildTerraformChart(t *testing.T) {
	chart := wavefront.Chart{
		Name:        "test_chart",
		Sources:     []wavefront.Source{},
		Description: "A chart",
		Units:       "unit",
	}

	result := buildTerraformChart(chart)
	if result["name"] != "test_chart" {
		t.Errorf("expected test_chart, got %s", result["name"])
	}
	if len(result["source"].([]map[string]interface{})) != 0 {
		t.Errorf("Expected empty array, got Array of lenth %d", len(result["source"].([]map[string]interface{})))
	}
	if result["description"] != "A chart" {
		t.Errorf("expected test_chart, got %s", result["description"])
	}
	if result["units"] != "unit" {
		t.Errorf("expected test_chart, got %s", result["units"])
	}

	chartWithSources := wavefront.Chart{
		Name: "I have sources",
		Sources: []wavefront.Source{
			{
				Name:  "chart 1",
				Query: "ts()",
			},
			{
				Name:  "chart 2",
				Query: "ts()",
			},
		},
	}

	resultWithCharts := buildTerraformChart(chartWithSources)
	if resultWithCharts["name"] != "I have sources" {
		t.Errorf("Expected \"I have sources\", got %s", resultWithCharts["name"])
	}
	if len(resultWithCharts["source"].([]map[string]interface{})) != 2 {
		t.Errorf("Expected empty array, got Array of lenth %d", len(result["source"].([]map[string]interface{})))
	}
}

func TestBuildTerraformChartSettings(t *testing.T) {
	chartSettings := wavefront.ChartSetting{
		Type:     "line",
		LineType: "linear",
	}
	result := buildTerraformChartSettings(chartSettings)
	if result["type"] != chartSettings.Type {
		t.Errorf("Expected %s, got %s", chartSettings.Type, result["type"])
	}
	if result["line_type"] != chartSettings.LineType {
		t.Errorf("Expected %s, got %s", chartSettings.LineType, result["line_type"])
	}
}

func TestBuildChartSettings(t *testing.T) {
	settings0 := make(map[string]interface{})
	settings0["type"] = "line"
	settings0["line_type"] = "linear"
	settings := []interface{}{
		settings0,
	}

	result := buildChartSettings(&settings)

	if result == nil {
		t.Errorf("Expected chart settings for %v", settings)
	}

	if result.Type != settings0["type"] {
		t.Errorf("Expected chart type %s, got %v", settings0["type"], result.Type)
	}

	if result.LineType != settings0["line_type"] {
		t.Errorf("Expected line type %s, got %v", settings0["line_type"], result.LineType)
	}
}

func TestBuildSections(t *testing.T) {
	section0 := make(map[string]interface{})
	section0["name"] = "section 0"
	section0["row"] = []interface{}{}

	section1 := make(map[string]interface{})
	section1["name"] = "section 1"
	section1["row"] = []interface{}{}

	sections := []interface{}{
		section0,
		section1,
	}

	result := buildSections(&sections)
	if len(*result) != 2 {
		t.Errorf("Expected 2 sections for %d", len(*result))
	}
}

func TestBuildRows(t *testing.T) {
	row0 := make(map[string]interface{})
	row0["chart"] = []interface{}{}

	row1 := make(map[string]interface{})
	row1["chart"] = []interface{}{}

	rows := []interface{}{
		row0,
		row1,
	}

	result := buildRows(&rows)
	if len(*result) != 2 {
		t.Errorf("Expected 2 rows for %d", len(*result))
	}
}

func TestBuildCharts(t *testing.T) {
	chart0 := make(map[string]interface{})
	chart0["name"] = "chart 0"
	chart0["description"] = "desc"
	chart0["units"] = "unit"
	chart0["source"] = []interface{}{}
	chart0["summarization"] = "MEAN"
	chart0["chart_setting"] = []interface{}{
		map[string]interface{}{
			"type": "linear",
		},
	}

	chart1 := make(map[string]interface{})
	chart1["name"] = "chart 1"
	chart1["description"] = "desc"
	chart1["units"] = "unit"
	chart1["source"] = []interface{}{}
	chart1["summarization"] = "MEAN"
	chart1["chart_setting"] = []interface{}{
		map[string]interface{}{
			"type": "linear",
		},
	}

	charts := []interface{}{
		chart0,
		chart1,
	}

	result := buildCharts(&charts)
	if len(*result) != 2 {
		t.Errorf("Expected 2 charts for %d", len(*result))
	}
	for i, r := range *result {
		if r.Name != fmt.Sprintf("chart %d", i) {
			t.Errorf("Expected chart %d for %s", i, r.Name)
		}
	}
}

func TestBuildSources(t *testing.T) {
	source0 := make(map[string]interface{})
	source0["name"] = "source 0"
	source0["query"] = "source 0"

	source1 := make(map[string]interface{})
	source1["name"] = "source 1"
	source1["query"] = "source 1"

	sources := []interface{}{
		source0,
		source1,
	}

	result := buildSources(&sources)
	if len(*result) != 2 {
		t.Errorf("Expected 2 sources for %d", len(*result))
	}
	for i, r := range *result {
		if r.Name != fmt.Sprintf("source %d", i) {
			t.Errorf("Expected source %d for %s", i, r.Name)
		}
	}
}

func TestBuildParameterDetails(t *testing.T) {
	param0 := make(map[string]interface{})
	param0["name"] = "source 0"
	param0["label"] = "source 0"
	param0["default_value"] = "test"
	param0["hide_from_view"] = true
	param0["parameter_type"] = "SIMPLE"
	param0["values_to_readable_strings"] = map[string]interface{}{
		"test": "test",
	}

	params := []interface{}{
		param0,
	}

	result := buildParameterDetails(&params)
	for k, v := range *result {
		if k != "source 0" {
			t.Errorf("Expected k 'source 0' got %s", k)
		}
		if v.Label != "source 0" {
			t.Errorf("Expected label 'source 0' got %s", k)
		}
	}

}

func TestAccWavefrontDashboard_Basic(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.test_dashboard", &record),
					testAccCheckWavefrontDashboardAttributes(&record),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "display_section_table_of_contents", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "display_query_parameters", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "name", "Terraform Test Dashboard"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "description", "testing, testing"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "url", "tftestcreate"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.name", "section 1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.row.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.row.0.chart.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.row.0.chart.0.description", "chart number 1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.row.0.chart.0.name", "chart 1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.row.0.chart.0.units", "something per unit"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.row.0.chart.0.source.0.name", "source name"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "section.0.row.0.chart.0.source.0.query", "ts()"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.0.default_value", "Label"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.0.hide_from_view", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.0.label", "param1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.0.name", "param1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.0.parameter_type", "SIMPLE"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.0.values_to_readable_strings.%", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "parameter_details.0.values_to_readable_strings.Label", "test"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_Updated(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.test_dashboard", &record),
					testAccCheckWavefrontDashboardAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "name", "Terraform Test Dashboard"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "display_section_table_of_contents", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "display_query_parameters", "true"),
				),
			},
			{
				Config: testAccCheckWavefrontDashboard_new_value(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.test_dashboard", &record),
					testAccCheckWavefrontDashboardAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "name", "Terraform Test Dashboard Updated"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "display_section_table_of_contents", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboard", "display_query_parameters", "false"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_Multiple(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_multiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.test_dashboarda", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboarda", "name", "Terraform Test Dashboard Multi A"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboardb", "name", "Terraform Test Dashboard Multi B"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.test_dashboardc", "name", "Terraform Test Dashboard Multi C"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_ListParam(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_ListParam(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.list_param_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.list_param_dash", "parameter_details.0.parameter_type", "LIST"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.list_param_dash", "parameter_details.0.values_to_readable_strings.%", "5"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.list_param_dash", "parameter_details.0.values_to_readable_strings.Label", "test"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.list_param_dash", "parameter_details.0.values_to_readable_strings.Label2", "test2"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_DynamicSourceParam(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_DynamicSourceParam(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.dynamic_source_param_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.parameter_type", "DYNAMIC"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.%", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.defaultQuery", "bubbles-fullstack-i-09ee71320cf9ba3c4.use1a.apelabs.net"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.dynamic_field_type", "SOURCE"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.query_value", "ts(servers.cpu-3.cpu-nice)"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_DynamicSourceTagParam(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_DynamicSourceTagParam(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.dynamic_source_param_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.parameter_type", "DYNAMIC"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.%", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.defaultQuery", "alpha"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.dynamic_field_type", "SOURCE_TAG"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.query_value", "ts(aws.lambda.invocations.average, tag=*)"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_DynamicMetricNameParam(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_DynamicMetricNameParam(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.dynamic_source_param_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.parameter_type", "DYNAMIC"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.%", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.defaultQuery", "aws.lambda.invocations.average"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.dynamic_field_type", "METRIC_NAME"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.query_value", "ts(aws.lambda.invocations.average)"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_DynamicTagKeyParam(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_DynamicTagKeyParam(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.dynamic_source_param_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.parameter_type", "DYNAMIC"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.%", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.defaultQuery", "aws.lambda.invocations.average"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.dynamic_field_type", "TAG_KEY"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.query_value", "ts(aws.lambda.invocations.average, test=test)"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.tag_key", "test"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_DynamicMatchingSourceTagParam(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_DynamicMatchingSourceTagParam(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.dynamic_source_param_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.parameter_type", "DYNAMIC"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.%", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.values_to_readable_strings.defaultQuery", "dev-elasticsearch"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.dynamic_field_type", "MATCHING_SOURCE_TAG"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.dynamic_source_param_dash", "parameter_details.0.query_value", "ts(aws.ec2.diskwritebytes.average)"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboard_Linear_ChartSettings(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_Linear_ChartSettings(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.chart_settings_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.name", "chart 1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.auto_column_tags", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.column_tags", "deprecated"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.custom_tags.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.expected_data_spacing", "120"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_display_stats.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_enabled", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_field", "CURRENT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_limit", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_sort", "TOP"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_hide_label", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_position", "RIGHT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_use_raw_stats", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.invert_dynamic_legend_hover_control", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.line_type", "linear"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.max", "100"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.min", "0"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.type", "line"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1max", "100"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_units", "units"),
				),
			},
		},
	})
}
func TestAccWavefrontDashboard_Table_ChartSettings(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_Table_ChartSettings(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.chart_settings_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.name", "chart 1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.auto_column_tags", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.column_tags", "deprecated"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.custom_tags.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.expected_data_spacing", "120"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_display_stats.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_enabled", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_field", "CURRENT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_limit", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_sort", "TOP"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_hide_label", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_position", "RIGHT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_use_raw_stats", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.group_by_source", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.line_type", "linear"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.num_tags", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.show_hosts", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.show_labels", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.show_raw_values", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sort_values_descending", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.tag_mode", "all"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.type", "table"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.windowing", "full"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.window_size", "10"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1max", "100"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_units", "units"),
				),
			},
		},
	})
}
func TestAccWavefrontDashboard_Sparkline_ChartSettings(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_Sparkline_ChartSettings(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.chart_settings_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.name", "chart 1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.auto_column_tags", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.column_tags", "deprecated"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.custom_tags.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.expected_data_spacing", "120"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_display_stats.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_enabled", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_field", "CURRENT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_limit", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_sort", "TOP"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_hide_label", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_position", "RIGHT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_use_raw_stats", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.line_type", "linear"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.max", "100"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.min", "0"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_decimal_precision", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_display_color", "rgba(1,1,1,1)"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_display_font_size", "14"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_display_horizontal_position", "LEFT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_display_postfix", "postfix"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_display_prefix", "prefix"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_display_value_type", "VALUE"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_display_vertical_position", "deprecated"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_fill_color", "rgba(1,1,1,1)"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_line_color", "rgba(1,1,1,1)"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_size", "BOTTOM"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_value_color_map_colors.#", "3"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_value_color_map_apply_to", "TEXT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_value_color_map_values.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.sparkline_value_text_map_text.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.type", "line"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1max", "100"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_units", "units"),
				),
			},
		},
	})
}
func TestAccWavefrontDashboard_Markdown_ChartSettings(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboard_Markdown_ChartSettings(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardExists("wavefront_dashboard.chart_settings_dash", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.name", "chart 1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.#", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.auto_column_tags", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.column_tags", "deprecated"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.custom_tags.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_display_stats.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_enabled", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_field", "CURRENT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_limit", "1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_filter_sort", "TOP"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_hide_label", "false"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_position", "RIGHT"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.fixed_legend_use_raw_stats", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.num_tags", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.plain_markdown_content", "markdown content"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.tag_mode", "all"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.type", "markdown-widget"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.windowing", "full"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.window_size", "10"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.max", "100"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.min", "0"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y0_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1max", "100"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_scale_si_by_1024", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_unit_autoscaling", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard.chart_settings_dash", "section.0.row.0.chart.0.chart_setting.0.y1_units", "units"),
				),
			},
		},
	})
}

func testAccCheckWavefrontDashboardDestroy(s *terraform.State) error {

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

func testAccCheckWavefrontDashboardAttributes(dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dashboard.Name != "Terraform Test Dashboard" {
			return fmt.Errorf("Bad value: %s", dashboard.Name)
		}

		return nil
	}
}

func testAccCheckWavefrontDashboardAttributesUpdated(dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dashboard.Name != "Terraform Test Dashboard Updated" {
			return fmt.Errorf("Bad value: %s", dashboard.Name)
		}

		return nil
	}
}

func testAccCheckWavefrontDashboardExists(n string, dashboard *wavefront.Dashboard) resource.TestCheckFunc {
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

func testAccCheckWavefrontDashboard_basic() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "test_dashboard" {
  name = "Terraform Test Dashboard"
  description = "testing, testing"
  url = "tftestcreate"
  display_section_table_of_contents = true
  display_query_parameters = true

  section{
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "SIMPLE"
	values_to_readable_strings = {
		Label = "test"
	}
  }
  tags = [
    "b",
    "terraform",
    "a",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_new_value() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "test_dashboard" {
  name = "Terraform Test Dashboard Updated"
  description = "testing, testing"
  url = "tftestcreate"
  display_section_table_of_contents = false
  display_query_parameters = false
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "SIMPLE"
    values_to_readable_strings = {
      Label = "test"
    }
  }
  tags = [
    "b",
    "terraform",
    "a",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_multiple() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "test_dashboarda" {
  name = "Terraform Test Dashboard Multi A"
  description = "testing, testing"
  url = "tftestmultia"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "SIMPLE"
    values_to_readable_strings = {
      Label = "test"
    }
  }
  tags = [
    "terraform",
    "test"
  ]
}
resource "wavefront_dashboard" "test_dashboardb" {
  name = "Terraform Test Dashboard Multi B"
  description = "testing, testing"
  url = "tftestmultib"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "SIMPLE"
    values_to_readable_strings = {
      Label = "test"
    }
  }
  tags = [
    "terraform",
    "test"
  ]
}
resource "wavefront_dashboard" "test_dashboardc" {
  name = "Terraform Test Dashboard Multi C"
  description = "testing, testing"
  url = "tftestmultic"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "SIMPLE"
    values_to_readable_strings = {
      Label = "test"
    }
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_ListParam() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "list_param_dash" {
  name = "Terraform Test Dashboard Updated"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "LIST"
    values_to_readable_strings = {
      b = "c",
      Label = "test",
      c = "d"
      Label2 = "test2",
      a = "b"
    }
  }
  parameter_details {
    name = "param2"
    label = "param2"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "LIST"
    values_to_readable_strings = {
      b = "c",
      Label = "test",
      c = "d"
      Label2 = "test2",
      a = "b"
    }
  }
  parameter_details {
    name = "param3"
    label = "param3"
    default_value = "Label"
    hide_from_view = false
    parameter_type = "LIST"
    values_to_readable_strings = {
      b = "c",
      Label = "test",
      c = "d"
      Label2 = "test2",
      a = "b"
    }
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_DynamicSourceParam() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "dynamic_source_param_dash" {
  name = "Terraform Dynamic Source Param"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "bubbles-fullstack-i-09ee71320cf9ba3c4.use1a.apelabs.net"
    }
    dynamic_field_type = "SOURCE"
    query_value = "ts(servers.cpu-3.cpu-nice)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_DynamicSourceTagParam() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "dynamic_source_param_dash" {
  name = "Terraform Dynamic Source Tag Param"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "alpha"
    }
    dynamic_field_type = "SOURCE_TAG"
    query_value = "ts(aws.lambda.invocations.average, tag=*)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_DynamicMetricNameParam() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "dynamic_source_param_dash" {
  name = "Terraform Dynamic Source Tag Param"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "aws.lambda.invocations.average"
    }
    dynamic_field_type = "METRIC_NAME"
    query_value = "ts(aws.lambda.invocations.average)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_DynamicTagKeyParam() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "dynamic_source_param_dash" {
  name = "Terraform Dynamic Source Tag Param"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "aws.lambda.invocations.average"
    }
    dynamic_field_type = "TAG_KEY"
    query_value = "ts(aws.lambda.invocations.average, test=test)"
    tag_key = "test"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_DynamicMatchingSourceTagParam() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "dynamic_source_param_dash" {
  name = "Terraform Dynamic Source Tag Param"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
          query = "ts()"
        }
        chart_setting {
          type = "linear"
        }
        summarization = "MEAN"
      }
    }
  }
  parameter_details {
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "dev-elasticsearch"
    }
    dynamic_field_type = "MATCHING_SOURCE_TAG"
    query_value = "ts(aws.ec2.diskwritebytes.average)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}

func testAccCheckWavefrontDashboard_Linear_ChartSettings() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "chart_settings_dash" {
  name = "Terraform Chart Settings"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
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
          invert_dynamic_legend_hover_control = true
          line_type = "linear"
          max = 100
          min = 0
          type = "line"
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
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "dev-elasticsearch"
    }
    dynamic_field_type = "MATCHING_SOURCE_TAG"
    query_value = "ts(aws.ec2.diskwritebytes.average)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}
func testAccCheckWavefrontDashboard_Table_ChartSettings() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "chart_settings_dash" {
  name = "Terraform Chart Settings"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
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
          line_type = "linear"
          num_tags = 2
          show_hosts = false
          show_labels = false
          show_raw_values = false
          sort_values_descending = true
          tag_mode = "all"
          type = "table"
          windowing = "full"
          window_size = 10
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
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "dev-elasticsearch"
    }
    dynamic_field_type = "MATCHING_SOURCE_TAG"
    query_value = "ts(aws.ec2.diskwritebytes.average)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}
func testAccCheckWavefrontDashboard_Sparkline_ChartSettings() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "chart_settings_dash" {
  name = "Terraform Chart Settings"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
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
          line_type = "linear"
          max = 100
          min = 0
          sparkline_decimal_precision = 1
          sparkline_display_color = "rgba(1,1,1,1)"
          sparkline_display_font_size = 14
          sparkline_display_horizontal_position = "LEFT"
          sparkline_display_postfix = "postfix"
          sparkline_display_prefix = "prefix"
          sparkline_display_value_type = "VALUE"
          sparkline_display_vertical_position = "deprecated"
          sparkline_fill_color = "rgba(1,1,1,1)"
          sparkline_line_color = "rgba(1,1,1,1)"
          sparkline_size = "BOTTOM"
          sparkline_value_color_map_apply_to = "TEXT"
          sparkline_value_color_map_colors = ["rgba(1,1,1,1)", "rgba(2,2,2,2)", "rgba(3,3,3,3)"]
          sparkline_value_color_map_values = [ 1, 2 ]
          sparkline_value_text_map_text = ["a"]
          sparkline_value_text_map_thresholds = [1]
          type = "line"
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
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "dev-elasticsearch"
    }
    dynamic_field_type = "MATCHING_SOURCE_TAG"
    query_value = "ts(aws.ec2.diskwritebytes.average)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}
func testAccCheckWavefrontDashboard_Markdown_ChartSettings() string {
	return fmt.Sprintf(`
resource "wavefront_dashboard" "chart_settings_dash" {
  name = "Terraform Chart Settings"
  description = "testing, testing"
  url = "tftestcreate"
  section {
    name = "section 1"
    row {
      chart {
        name = "chart 1"
        description = "chart number 1"
        units = "something per unit"
        source {
          name = "source name"
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
    name = "param1"
    label = "param1"
    default_value = "defaultQuery"
    hide_from_view = false
    parameter_type = "DYNAMIC"
    values_to_readable_strings = {
      defaultQuery = "dev-elasticsearch"
    }
    dynamic_field_type = "MATCHING_SOURCE_TAG"
    query_value = "ts(aws.ec2.diskwritebytes.average)"
  }
  tags = [
    "terraform",
    "test"
  ]
}
`)
}
