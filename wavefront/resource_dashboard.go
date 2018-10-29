package wavefront_plugin

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spaceapegames/go-wavefront"
	"sort"
	"strings"
)

// Terraform Resource Declaration
func resourceDashboard() *schema.Resource {

	source := &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "A collection of Sources for a Chart",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the Source",
				},
				"query": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Query for the Source",
				},
				"disabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether to disabled the source from being displayed",
				},
				"scatter_plot_source": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "Y",
				},
				"query_builder_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether the query builder should be enabled",
				},
				"source_description": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Description of the source",
				},
			},
		},
	}

	chartSettings := &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_column_tags": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"column_tags": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"custom_tags": {
					Type:     schema.TypeSet,
					Optional: true,
				},
				"expected_data_spacing": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"fixed_legend_display_stats": {
					Type:     schema.TypeSet,
					Optional: true,
				},
				"fixed_legend_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"fixed_legend_filter_field": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"fixed_legend_filter_limit": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"fixed_legend_filter_sort": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"fixed_legend_hide_label": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"fixed_legend_position": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"fixed_legend_use_raw_stats": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"group_by_source": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"invert_dynamic_legend_hover_control": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"line_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"max": {
					Type:     schema.TypeFloat,
					Optional: true,
					Default:  0.0,
				},
				"min": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"num_tags": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"plain_markdown_content": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"show_hosts": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"show_labels": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"show_raw_values": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"sort_values_descending": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"sparkline_decimal_precision": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"sparkline_display_color": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_display_font_size": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_display_horizontal_position": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_display_postfix": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_display_prefix": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_display_value_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_display_vertical_position": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_fill_color": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_line_color": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_size": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_value_color_map_apply_to": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sparkline_value_color_map_colors": {
					Type:     schema.TypeSet,
					Optional: true,
				},
				"sparkline_value_color_map_values": {
					Type:     schema.TypeSet,
					Optional: true,
				},
				"sparkline_value_color_map_values_v2": {
					Type:     schema.TypeSet,
					Optional: true,
				},
				"sparkline_value_text_map_text": {
					Type:     schema.TypeSet,
					Optional: true,
				},
				"sparkline_value_text_map_thresholds": {
					Type:     schema.TypeSet,
					Optional: true,
				},
				"stack_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"tag_mode": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"time_based_coloring": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"windowing": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"window_size": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"x_max": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"x_min": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"y0_scale_sib_y1024": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"y0_unit_autoscaling": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"y1_max": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"y1_min": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"y1_scale_sib_y1024": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"y1_unit_autoscaling": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"y1_units": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"y_max": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"y_min": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
			},
		},
	}

	chart := &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "A collection of chart",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the Chart",
				},
				"description": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Description of the chart",
				},
				"units": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Units of measurements for the chart",
				},
				"source": source,
				"summarization": {
					Type:        schema.TypeString,
					Default:     "MEAN",
					Description: "Strategy to use when aggregating metric point [MEAN, MEDIAN, MIN, MAX, COUNT, SUM, LAST, FIRST]",
					Optional:    true,
				},
				"base": {
					Type:        schema.TypeInt,
					Default:     0,
					Description: "unknown usage, defaults to 0",
					Optional:    true,
				},
				"include_obsolete_metrics": {
					Type:        schema.TypeBool,
					Description: "Include obsolete metrics older than 4 weeks ago into current time window",
					Optional:    true,
				},
				"interpolate_points": {
					Type:        schema.TypeBool,
					Description: "Interpolate points that existed in past/future into current time window",
					Optional:    true,
				},
				"no_default_events": {
					Type:        schema.TypeBool,
					Description: "Don't include default events on the chart",
					Optional:    true,
				},
				"chart_settings": chartSettings,
			},
		},
	}

	row := &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "Rows containing chart. Rows belong in Sections",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"chart": chart,
			},
		},
	}

	section := &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "Sections of a Dashboard",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the Sections",
				},
				"row": row,
			},
		},
	}

	parameterDetail := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "",
					Required:    true,
				},
				"label": {
					Type:     schema.TypeString,
					Required: true,
				},
				"default_value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"hide_from_view": {
					Type:     schema.TypeBool,
					Required: true,
				},
				"parameter_type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"values_to_readable_strings": {
					Type:        schema.TypeMap,
					Required:    true,
					Description: "Map of [string]string. At least one of the keys must match the value of default_value.",
				},
				"query_value": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"tag_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"dynamic_field_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	return &schema.Resource{
		Create: resourceDashboardCreate,
		Read:   resourceDashboardRead,
		Update: resourceDashboardUpdate,
		Delete: resourceDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"section":           section,
			"parameter_details": parameterDetail,
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"chart_title_bg_color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"chart_title_color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"chart_title_scalar": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"default_end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"default_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"default_time_window": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_description": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"display_query_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"display_section_table_of_contents": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"event_filter_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "BYCHART",
				Description: "one of [NONE, ALL, BYCHART, AUTOMATIC, BYDASHBOARD, BYCHARTANDDASHBOARD]. Default: BYCHART",
			},
			"event_query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"favorite": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

// Construct a Terraform ParameterDetail
func buildTerraformParameterDetail(wavefrontParamDetail wavefront.ParameterDetail, name string) map[string]interface{} {
	parameterDetail := map[string]interface{}{}

	parameterDetail["name"] = name
	parameterDetail["label"] = wavefrontParamDetail.Label
	parameterDetail["parameter_type"] = wavefrontParamDetail.ParameterType
	parameterDetail["hide_from_view"] = wavefrontParamDetail.HideFromView
	parameterDetail["default_value"] = wavefrontParamDetail.DefaultValue
	parameterDetail["values_to_readable_strings"] = wavefrontParamDetail.ValuesToReadableStrings
	parameterDetail["query_value"] = wavefrontParamDetail.QueryValue
	parameterDetail["tag_key"] = wavefrontParamDetail.TagKey
	parameterDetail["dynamic_field_type"] = wavefrontParamDetail.DynamicFieldType
	return parameterDetail
}

// Construct a Terraform Section
func buildTerraformSection(wavefrontSection wavefront.Section) map[string]interface{} {
	section := map[string]interface{}{}
	section["name"] = wavefrontSection.Name
	rows := []map[string]interface{}{}
	for _, wavefrontRow := range wavefrontSection.Rows {
		rows = append(rows, buildTerraformRow(wavefrontRow))
	}
	section["row"] = rows

	return section
}

// Construct a Wavefront Row
func buildTerraformRow(wavefrontRow wavefront.Row) map[string]interface{} {
	row := map[string]interface{}{}

	charts := []map[string]interface{}{}
	for _, wavefrontRow := range wavefrontRow.Charts {
		charts = append(charts, buildTerraformChart(wavefrontRow))
	}
	row["chart"] = charts

	return row
}

// Construct a Wavefront Chart
func buildTerraformChart(wavefrontChart wavefront.Chart) map[string]interface{} {
	chart := map[string]interface{}{}
	chart["name"] = wavefrontChart.Name
	chart["description"] = wavefrontChart.Description

	chart["units"] = wavefrontChart.Units
	sources := []map[string]interface{}{}
	for _, wavefrontSource := range wavefrontChart.Sources {
		sources = append(sources, buildTerraformSource(wavefrontSource))
	}
	chart["source"] = sources
	chart["summarization"] = wavefrontChart.Summarization
	chart["include_obsolete_metrics"] = wavefrontChart.IncludeObsoleteMetrics
	chart["interpolate_points"] = wavefrontChart.InterpolatePoints
	chart["no_default_events"] = wavefrontChart.NoDefaultEvents

	chart["chart_settings"] = buildTerraformChartSetting(wavefrontChart.ChartSettings)

	return chart
}

// Construct a Wavefront Source
func buildTerraformSource(wavefrontSource wavefront.Source) map[string]interface{} {
	source := map[string]interface{}{}
	source["name"] = wavefrontSource.Name
	source["query"] = wavefrontSource.Query
	source["disabled"] = wavefrontSource.Disabled
	source["scatter_plot_source"] = wavefrontSource.ScatterPlotSource
	source["query_builder_enabled"] = wavefrontSource.QuerybuilderEnabled
	source["source_description"] = wavefrontSource.SourceDescription

	return source
}

func buildTerraformChartSetting(wavefrontChartSetting *wavefront.ChartSetting) map[string]interface{} {
	chartSetting := map[string]interface{}{}

	if wavefrontChartSetting == nil {
		return chartSetting
	}
	fmt.Println("buildTerraformChartSetting", spew.Sdump(wavefrontChartSetting))
	chartSetting["auto_column_tags"] = wavefrontChartSetting.AutoColumnTags
	chartSetting["column_tags"] = wavefrontChartSetting.ColumnTags
	chartSetting["custom_tags"] = wavefrontChartSetting.CustomTags
	chartSetting["expected_data_spacing"] = wavefrontChartSetting.ExpectedDataSpacing
	chartSetting["fixed_legend_display_stats"] = wavefrontChartSetting.FixedLegendDisplayStats
	chartSetting["fixed_legend_enabled"] = wavefrontChartSetting.FixedLegendEnabled
	chartSetting["fixed_legend_filter_field"] = wavefrontChartSetting.FixedLegendFilterField
	chartSetting["fixed_legend_filter_limit"] = wavefrontChartSetting.FixedLegendFilterLimit
	chartSetting["fixed_legend_filter_sort"] = wavefrontChartSetting.FixedLegendFilterSort
	chartSetting["fixed_legend_hide_label"] = wavefrontChartSetting.FixedLegendHideLabel
	chartSetting["fixed_legend_position"] = wavefrontChartSetting.FixedLegendPosition
	chartSetting["fixed_legend_use_raw_stats"] = wavefrontChartSetting.FixedLegendUseRawStats
	chartSetting["group_by_source"] = wavefrontChartSetting.GroupBySource
	chartSetting["invert_dynamic_legend_hover_control"] = wavefrontChartSetting.InvertDynamicLegendHoverControl
	chartSetting["line_type"] = wavefrontChartSetting.LineType
	chartSetting["max"] = wavefrontChartSetting.Max
	chartSetting["min"] = wavefrontChartSetting.Min
	chartSetting["num_tags"] = wavefrontChartSetting.NumTags
	chartSetting["plain_markdown_content"] = wavefrontChartSetting.PlainMarkdownContent
	chartSetting["show_hosts"] = wavefrontChartSetting.ShowHosts
	chartSetting["show_labels"] = wavefrontChartSetting.ShowLabels
	chartSetting["show_raw_values"] = wavefrontChartSetting.ShowRawValues
	chartSetting["sort_values_descending"] = wavefrontChartSetting.SortValuesDescending
	chartSetting["sparkline_decimal_precision"] = wavefrontChartSetting.SparklineDecimalPrecision
	chartSetting["sparkline_display_color"] = wavefrontChartSetting.SparklineDisplayColor
	chartSetting["sparkline_display_font_size"] = wavefrontChartSetting.SparklineDisplayFontSize
	chartSetting["sparkline_display_horizontal_position"] = wavefrontChartSetting.SparklineDisplayHorizontalPosition
	chartSetting["sparkline_display_postfix"] = wavefrontChartSetting.SparklineDisplayPostfix
	chartSetting["sparkline_display_prefix"] = wavefrontChartSetting.SparklineDisplayPrefix
	chartSetting["sparkline_display_value_type"] = wavefrontChartSetting.SparklineDisplayValueType
	chartSetting["sparkline_display_vertical_position"] = wavefrontChartSetting.SparklineDisplayVerticalPosition
	chartSetting["sparkline_fill_color"] = wavefrontChartSetting.SparklineFillColor
	chartSetting["sparkline_line_color"] = wavefrontChartSetting.SparklineLineColor
	chartSetting["sparkline_size"] = wavefrontChartSetting.SparklineSize
	chartSetting["sparkline_value_color_map_apply_to"] = wavefrontChartSetting.SparklineValueColorMapApplyTo
	chartSetting["sparkline_value_color_map_colors"] = wavefrontChartSetting.SparklineValueColorMapColors
	chartSetting["sparkline_value_color_map_values"] = wavefrontChartSetting.SparklineValueColorMapValues
	chartSetting["sparkline_value_color_map_values_v2"] = wavefrontChartSetting.SparklineValueColorMapValuesV2
	chartSetting["sparkline_value_text_map_text"] = wavefrontChartSetting.SparklineValueTextMapText
	chartSetting["sparkline_value_text_map_thresholds"] = wavefrontChartSetting.SparklineValueTextMapThresholds
	chartSetting["stack_type"] = wavefrontChartSetting.StackType
	chartSetting["tag_mode"] = wavefrontChartSetting.TagMode
	chartSetting["time_based_coloring"] = wavefrontChartSetting.TimeBasedColoring
	chartSetting["type"] = wavefrontChartSetting.Type
	chartSetting["windowing"] = wavefrontChartSetting.Windowing
	chartSetting["window_size"] = wavefrontChartSetting.WindowSize
	chartSetting["x_max"] = wavefrontChartSetting.Xmax
	chartSetting["x_min"] = wavefrontChartSetting.Xmin
	chartSetting["y0_scale_sib_y1024"] = wavefrontChartSetting.Y0ScaleSIBy1024
	chartSetting["y0_unit_autoscaling"] = wavefrontChartSetting.Y0UnitAutoscaling
	chartSetting["y1_max"] = wavefrontChartSetting.Y1Max
	chartSetting["y1_min"] = wavefrontChartSetting.Y1Min
	chartSetting["y1_scale_sib_y1024"] = wavefrontChartSetting.Y1ScaleSIBy1024
	chartSetting["y1_unit_autoscaling"] = wavefrontChartSetting.Y1UnitAutoscaling
	chartSetting["y1_units"] = wavefrontChartSetting.Y1Units
	chartSetting["y_max"] = wavefrontChartSetting.Ymax
	chartSetting["y_min"] = wavefrontChartSetting.Ymin

	return chartSetting
}

// Construct a Wavefront Section
func buildSections(terraformSections *[]interface{}) *[]wavefront.Section {
	wavefrontSections := make([]wavefront.Section, len(*terraformSections))

	for i, t := range *terraformSections {
		t := t.(map[string]interface{})

		terraformRows := t["row"].([]interface{})

		wavefrontSections[i] = wavefront.Section{
			Name: t["name"].(string),
			Rows: *buildRows(&terraformRows),
		}
	}
	return &wavefrontSections
}

// Construct a Wavefront Row
func buildRows(terraformRows *[]interface{}) *[]wavefront.Row {
	wavefrontRows := make([]wavefront.Row, len(*terraformRows))

	for i, t := range *terraformRows {
		t := t.(map[string]interface{})

		terraformCharts := t["chart"].([]interface{})

		wavefrontRows[i] = wavefront.Row{
			Charts: *buildCharts(&terraformCharts),
		}
	}

	return &wavefrontRows
}

// Construct a Wavefront Chart
func buildCharts(terraformCharts *[]interface{}) *[]wavefront.Chart {
	wavefrontCharts := make([]wavefront.Chart, len(*terraformCharts))

	for i, t := range *terraformCharts {
		t := t.(map[string]interface{})

		terraformSources := t["source"].([]interface{})
		terraformChartSettings := map[string]interface{}{}
		chartSettingsNotNil := false
		if t["chart_settings"] != nil {
			chartSettingsNotNil = true
			terraformChartSettings = t["chart_settings"].(map[string]interface{})
		}

		wavefrontCharts[i] = wavefront.Chart{
			Name:                   t["name"].(string),
			Sources:                *buildSources(&terraformSources),
			Description:            t["description"].(string),
			Units:                  t["units"].(string),
			Summarization:          t["summarization"].(string),
			IncludeObsoleteMetrics: t["include_obsolete_metrics"].(bool),
			InterpolatePoints:      t["interpolate_points"].(bool),
		}
		if t["no_default_events"] != nil {
			wavefrontCharts[i].NoDefaultEvents = t["no_default_events"].(bool)
		}
		if chartSettingsNotNil {
			wavefrontCharts[i].ChartSettings = buildchartSettings(terraformChartSettings)
		}
	}

	return &wavefrontCharts
}

func buildchartSettings(chartSettings interface{}) *wavefront.ChartSetting {
	wavefrontChartSetting := wavefront.ChartSetting{}

	notNil := false

	terraformChartSettings := chartSettings.(map[string]interface{})

	if terraformChartSettings == nil {
		notNil = true
		return &wavefrontChartSetting
	}
	fmt.Println("buildchartSettings", spew.Sdump(terraformChartSettings))
	if terraformChartSettings["auto_column_tags"] != nil {
		notNil = true
		wavefrontChartSetting.AutoColumnTags = terraformChartSettings["auto_column_tags"].(bool)
	}
	if terraformChartSettings["column_tags"] != nil {
		notNil = true
		wavefrontChartSetting.ColumnTags = terraformChartSettings["column_tags"].(string)
	}
	if terraformChartSettings["custom_tags"] != nil {
		notNil = true
		wavefrontChartSetting.ColumnTags = terraformChartSettings["custom_tags"].(string)
	}
	if terraformChartSettings["expected_data_spacing"] != nil {
		notNil = true
		wavefrontChartSetting.ExpectedDataSpacing = terraformChartSettings["expected_data_spacing"].(int)
	}
	if terraformChartSettings["fixed_legend_display_stats"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendDisplayStats = terraformChartSettings["fixed_legend_display_stats"].([]string)
	}
	if terraformChartSettings["fixed_legend_enabled"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendEnabled = terraformChartSettings["fixed_legend_enabled"].(bool)
	}
	if terraformChartSettings["fixed_legend_filter_field"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendFilterField = terraformChartSettings["fixed_legend_filter_field"].(string)
	}
	if terraformChartSettings["fixed_legend_filter_limit"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendFilterLimit = terraformChartSettings["fixed_legend_filter_limit"].(int)
	}
	if terraformChartSettings["fixed_legend_filter_sort"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendFilterSort = terraformChartSettings["fixed_legend_filter_sort"].(string)
	}
	if terraformChartSettings["fixed_legend_hide_label"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendHideLabel = terraformChartSettings["fixed_legend_hide_label"].(bool)
	}
	if terraformChartSettings["fixed_legend_position"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendPosition = terraformChartSettings["fixed_legend_position"].(string)
	}
	if terraformChartSettings["fixed_legend_use_raw_stats"] != nil {
		notNil = true
		wavefrontChartSetting.FixedLegendUseRawStats = terraformChartSettings["fixed_legend_use_raw_stats"].(bool)
	}
	if terraformChartSettings["group_by_source"] != nil {
		notNil = true
		wavefrontChartSetting.GroupBySource = terraformChartSettings["group_by_source"].(bool)
	}
	if terraformChartSettings["invert_dynamic_legend_hover_control"] != nil {
		notNil = true
		wavefrontChartSetting.InvertDynamicLegendHoverControl = terraformChartSettings["invert_dynamic_legend_hover_control"].(bool)
	}
	if terraformChartSettings["line_type"] != nil {
		notNil = true
		wavefrontChartSetting.LineType = terraformChartSettings["line_type"].(string)
	}
	if terraformChartSettings["max"] != nil {
		notNil = true
		wavefrontChartSetting.Max = terraformChartSettings["max"].(float32)
	}
	if terraformChartSettings["min"] != nil {
		notNil = true
		wavefrontChartSetting.Min = terraformChartSettings["min"].(float32)
	}
	if terraformChartSettings["num_tags"] != nil {
		notNil = true
		wavefrontChartSetting.NumTags = terraformChartSettings["num_tags"].(int)
	}
	if terraformChartSettings["plain_markdown_content"] != nil {
		notNil = true
		wavefrontChartSetting.PlainMarkdownContent = terraformChartSettings["plain_markdown_content"].(string)
	}
	if terraformChartSettings["show_hosts"] != nil {
		notNil = true
		wavefrontChartSetting.ShowHosts = terraformChartSettings["show_hosts"].(bool)
	}
	if terraformChartSettings["show_labels"] != nil {
		notNil = true
		wavefrontChartSetting.ShowLabels = terraformChartSettings["show_labels"].(bool)
	}
	if terraformChartSettings["show_raw_values"] != nil {
		notNil = true
		wavefrontChartSetting.ShowRawValues = terraformChartSettings["show_raw_values"].(bool)
	}
	if terraformChartSettings["sort_values_descending"] != nil {
		notNil = true
		wavefrontChartSetting.SortValuesDescending = terraformChartSettings["sort_values_descending"].(bool)
	}
	if terraformChartSettings["sparkline_decimal_precision"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDecimalPrecision = terraformChartSettings["sparkline_decimal_precision"].(int)
	}
	if terraformChartSettings["sparkline_display_color"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDisplayColor = terraformChartSettings["sparkline_display_color"].(string)
	}
	if terraformChartSettings["sparkline_display_font_size"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDisplayFontSize = terraformChartSettings["sparkline_display_font_size"].(string)
	}
	if terraformChartSettings["sparkline_display_horizontal_position"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDisplayHorizontalPosition = terraformChartSettings["sparkline_display_horizontal_position"].(string)
	}
	if terraformChartSettings["sparkline_display_postfix"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDisplayPostfix = terraformChartSettings["sparkline_display_postfix"].(string)
	}
	if terraformChartSettings["sparkline_display_prefix"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDisplayPrefix = terraformChartSettings["sparkline_display_prefix"].(string)
	}
	if terraformChartSettings["sparkline_display_value_type"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDisplayValueType = terraformChartSettings["sparkline_display_value_type"].(string)
	}
	if terraformChartSettings["sparkline_display_vertical_position"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineDisplayVerticalPosition = terraformChartSettings["sparkline_display_vertical_position"].(string)
	}
	if terraformChartSettings["sparkline_fill_color"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineFillColor = terraformChartSettings["sparkline_fill_color"].(string)
	}
	if terraformChartSettings["sparkline_line_color"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineLineColor = terraformChartSettings["sparkline_line_color"].(string)
	}
	if terraformChartSettings["sparkline_size"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineSize = terraformChartSettings["sparkline_size"].(string)
	}
	if terraformChartSettings["sparkline_value_color_map_apply_to"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineValueColorMapApplyTo = terraformChartSettings["sparkline_value_color_map_apply_to"].(string)
	}
	if terraformChartSettings["sparkline_value_color_map_colors"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineValueColorMapColors = terraformChartSettings["sparkline_value_color_map_colors"].([]string)
	}
	if terraformChartSettings["sparkline_value_color_map_values"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineValueColorMapValues = terraformChartSettings["sparkline_value_color_map_values"].([]int)
	}
	if terraformChartSettings["sparkline_value_color_map_values_v2"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineValueColorMapValuesV2 = terraformChartSettings["sparkline_value_color_map_values_v2"].([]int)
	}
	if terraformChartSettings["sparkline_value_text_map_text"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineValueTextMapText = terraformChartSettings["sparkline_value_text_map_text"].([]string)
	}
	if terraformChartSettings["sparkline_value_text_map_thresholds"] != nil {
		notNil = true
		wavefrontChartSetting.SparklineValueTextMapThresholds = terraformChartSettings["sparkline_value_text_map_thresholds"].([]int)
	}
	if terraformChartSettings["stack_type"] != nil {
		notNil = true
		wavefrontChartSetting.StackType = terraformChartSettings["stack_type"].(string)
	}
	if terraformChartSettings["tag_mode"] != nil {
		notNil = true
		wavefrontChartSetting.TagMode = terraformChartSettings["tag_mode"].(string)
	}
	if terraformChartSettings["time_based_coloring"] != nil {
		notNil = true
		wavefrontChartSetting.TimeBasedColoring = terraformChartSettings["time_based_coloring"].(bool)
	}
	if terraformChartSettings["type"] != nil {
		notNil = true
		wavefrontChartSetting.Type = terraformChartSettings["type"].(string)
	}
	if terraformChartSettings["windowing"] != nil {
		notNil = true
		wavefrontChartSetting.Windowing = terraformChartSettings["windowing"].(string)
	}
	if terraformChartSettings["window_size"] != nil {
		notNil = true
		wavefrontChartSetting.WindowSize = terraformChartSettings["window_size"].(int)
	}
	if terraformChartSettings["x_max"] != nil {
		notNil = true
		wavefrontChartSetting.Xmax = terraformChartSettings["x_max"].(float32)
	}
	if terraformChartSettings["sparkline_display_font_size"] != nil {
		notNil = true
		wavefrontChartSetting.Xmin = terraformChartSettings["x_min"].(float32)
	}
	if terraformChartSettings["y0_scale_sib_y1024"] != nil {
		notNil = true
		wavefrontChartSetting.Y0ScaleSIBy1024 = terraformChartSettings["y0_scale_sib_y1024"].(bool)
	}
	if terraformChartSettings["y0_unit_autoscaling"] != nil {
		notNil = true
		wavefrontChartSetting.Y0UnitAutoscaling = terraformChartSettings["y0_unit_autoscaling"].(bool)
	}
	if terraformChartSettings["y1_max"] != nil {
		notNil = true
		wavefrontChartSetting.Y1Max = terraformChartSettings["y1_max"].(float32)
	}
	if terraformChartSettings["y1_min"] != nil {
		notNil = true
		wavefrontChartSetting.Y1Min = terraformChartSettings["y1_min"].(float32)
	}
	if terraformChartSettings["y1_scale_sib_y1024"] != nil {
		notNil = true
		wavefrontChartSetting.Y1ScaleSIBy1024 = terraformChartSettings["y1_scale_sib_y1024"].(bool)
	}
	if terraformChartSettings["y1_unit_autoscaling"] != nil {
		notNil = true
		wavefrontChartSetting.Y1UnitAutoscaling = terraformChartSettings["y1_unit_autoscaling"].(bool)
	}
	if terraformChartSettings["y1_units"] != nil {
		notNil = true
		wavefrontChartSetting.Y1Units = terraformChartSettings["y1_units"].(string)
	}
	if terraformChartSettings["y_max"] != nil {
		notNil = true
		wavefrontChartSetting.Ymax = terraformChartSettings["y_max"].(float32)
	}
	if terraformChartSettings["y_min"] != nil {
		notNil = true
		wavefrontChartSetting.Ymin = terraformChartSettings["y_min"].(float32)
	}

	if !notNil {
		return nil
	}

	return &wavefrontChartSetting
}

// Construct a Wavefront Source
func buildSources(terraformSources *[]interface{}) *[]wavefront.Source {
	wavefrontSources := make([]wavefront.Source, len(*terraformSources))

	for i, t := range *terraformSources {
		t := t.(map[string]interface{})

		wavefrontSources[i] = wavefront.Source{
			Name:  t["name"].(string),
			Query: t["query"].(string),
		}
		if t["disabled"] != nil {
			wavefrontSources[i].Disabled = t["disabled"].(bool)
		}
		if t["scatter_plot_source"] != nil {
			wavefrontSources[i].ScatterPlotSource = t["scatter_plot_source"].(string)
		}
		if t["query_builder_enabled"] != nil {
			wavefrontSources[i].QuerybuilderEnabled = t["query_builder_enabled"].(bool)
		}
		if t["source_description"] != nil {
			wavefrontSources[i].SourceDescription = t["source_description"].(string)
		}
	}

	return &wavefrontSources
}

// Construct a Wavefront ParameterDetail
func buildParameterDetails(terraformParams *[]interface{}) *map[string]wavefront.ParameterDetail {
	wavefrontParams := map[string]wavefront.ParameterDetail{}

	for _, t := range *terraformParams {
		t := t.(map[string]interface{})
		name := t["name"].(string)
		valuesToReadableStrings := t["values_to_readable_strings"].(map[string]interface{})
		readableStrings := map[string]string{}

		for k, v := range valuesToReadableStrings {
			readableStrings[k] = v.(string)
		}

		wfParam := wavefront.ParameterDetail{
			Label:                   t["label"].(string),
			DefaultValue:            t["default_value"].(string),
			HideFromView:            t["hide_from_view"].(bool),
			ParameterType:           t["parameter_type"].(string),
			ValuesToReadableStrings: readableStrings,
		}
		if t["query_value"] != nil {
			wfParam.QueryValue = t["query_value"].(string)
		}
		if t["tag_key"] != nil {
			wfParam.TagKey = t["tag_key"].(string)
		}
		if t["dynamic_field_type"] != nil {
			wfParam.DynamicFieldType = t["dynamic_field_type"].(string)
		}

		wavefrontParams[name] = wfParam
	}

	return &wavefrontParams
}

// Construct a Wavefront Dashboard
func buildDashboard(d *schema.ResourceData) (*wavefront.Dashboard, error) {

	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	terraformSections := d.Get("section").([]interface{})
	terraformParams := d.Get("parameter_details").([]interface{})

	dash := &wavefront.Dashboard{
		Name:             d.Get("name").(string),
		ID:               d.Get("url").(string),
		Tags:             tags,
		Description:      d.Get("description").(string),
		Url:              d.Get("url").(string),
		Sections:         *buildSections(&terraformSections),
		ParameterDetails: *buildParameterDetails(&terraformParams),
		EventFilterType:  d.Get("event_filter_type").(string),
	}

	if d.Get("chart_title_bg_color") != nil {
		dash.ChartTitleBgColor = d.Get("chart_title_bg_color").(string)
	}
	if d.Get("chart_title_color") != nil {
		dash.ChartTitleColor = d.Get("chart_title_color").(string)
	}
	if d.Get("chart_title_scalar") != nil {
		dash.ChartTitleScalar = d.Get("chart_title_scalar").(int)
	}
	if d.Get("default_end_time") != nil {
		dash.DefaultEndTime = d.Get("default_end_time").(int)
	}
	if d.Get("default_start_time") != nil {
		dash.DefaultStartTime = d.Get("default_start_time").(int)
	}
	if d.Get("default_time_window") != nil {
		dash.DefaultTimeWindow = d.Get("default_time_window").(string)
	}
	if d.Get("display_description") != nil {
		dash.DisplayDescription = d.Get("display_description").(bool)
	}
	if d.Get("display_query_parameters") != nil {
		dash.DisplayQueryParameters = d.Get("display_query_parameters").(bool)
	}

	return dash, nil
}

// Create a Terraform Dashboard
func resourceDashboardCreate(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()
	fmt.Println(spew.Sdump(d))
	dashboard, err := buildDashboard(d)
	if err != nil {
		return fmt.Errorf("failed to parse dashboard, %s", err)
	}

	err = dashboards.Create(dashboard)
	if err != nil {
		return fmt.Errorf("failed to create dashboard, %s", err)
	}
	d.SetId(dashboard.ID)

	return nil
}

type Params []map[string]interface{}

func (p Params) Len() int      { return len(p) }
func (p Params) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p Params) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{p[i]["name"].(string), p[j]["name"].(string)})
}

// Read a Wavefront Dashboard
func resourceDashboardRead(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()
	dash := wavefront.Dashboard{
		ID: d.Id(),
	}

	// search for an dashboard with our id. We should receive 1 (Exact Match) or 0 (No Match)
	err := dashboards.Get(&dash)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Wavefront Dashboard %s. %s", d.Id(), err)
		}
	}

	// Use the Wavefront url as the Terraform ID
	d.SetId(dash.ID)
	d.Set("name", dash.Name)
	d.Set("description", dash.Description)
	d.Set("url", dash.Url)
	d.Set("event_filter_type", dash.EventFilterType)

	if dash.ChartTitleBgColor != "" {
		d.Set("chart_title_bg_color", dash.ChartTitleBgColor)
	}
	if dash.ChartTitleColor != "" {
		d.Set("chart_title_color", dash.ChartTitleColor)
	}

	sections := []map[string]interface{}{}
	for _, wavefrontSection := range dash.Sections {
		sections = append(sections, buildTerraformSection(wavefrontSection))
	}
	d.Set("section", sections)

	parameterDetails := []map[string]interface{}{}

	for k, v := range dash.ParameterDetails {
		parameterDetails = append(parameterDetails, buildTerraformParameterDetail(v, k))
	}

	sort.Sort(Params(parameterDetails))

	d.Set("parameter_details", parameterDetails)
	d.Set("tags", dash.Tags)

	return nil
}

func resourceDashboardUpdate(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()

	a, err := buildDashboard(d)
	if err != nil {
		return fmt.Errorf("failed to parse dashboard, %s", err)
	}
	// Update the dashboard on Wavefront
	err = dashboards.Update(a)
	if err != nil {
		return fmt.Errorf("error Updating Dashboard %s. %s", d.Get("name"), err)
	}
	return nil
}

func resourceDashboardDelete(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()
	dash := wavefront.Dashboard{
		ID: d.Id(),
	}

	err := dashboards.Get(&dash)
	if err != nil {
		return fmt.Errorf("error finding Wavefront Dashboard %s. %s", d.Id(), err)
	}

	// Delete the Dashboard
	err = dashboards.Delete(&dash)
	if err != nil {
		return fmt.Errorf("failed to delete Dashboard %s. %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
