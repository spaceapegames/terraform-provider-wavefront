package wavefront_plugin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spaceapegames/go-wavefront"
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

	chartSetting := &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Required:    true,
		Description: "Chart settings. Defaults to line charts",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_column_tags": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "deprecated",
				},
				"column_tags": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "deprecated",
				},
				"custom_tags": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "For the tabular view, a list of point tags to display when using the custom tag display mode",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"expected_data_spacing": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Threshold (in seconds) for time delta between consecutive points in a series above which a dotted line will replace a solid line in line plots. Default: 60s",
				},
				"fixed_legend_display_stats": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "For a chart with a fixed legend, a list of statistics to display in the legend",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"fixed_legend_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether to enable a fixed tabular legend adjacent to the chart",
				},
				"fixed_legend_filter_field": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Statistic to use for determining whether a series is displayed on the fixed legend = ['CURRENT', 'MEAN', 'MEDIAN', 'SUM', 'MIN', 'MAX', 'COUNT']",
				},
				"fixed_legend_filter_limit": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Number of series to include in the fixed legend",
				},
				"fixed_legend_filter_sort": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Whether to display Top- or Bottom-ranked series in the fixed legend = ['TOP', 'BOTTOM']",
				},
				"fixed_legend_hide_label": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "deprecated",
				},
				"fixed_legend_position": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Where the fixed legend should be displayed with respect to the chart = ['RIGHT', 'TOP', 'LEFT', 'BOTTOM']",
				},
				"fixed_legend_use_raw_stats": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "If true, the legend uses non-summarized stats instead of summarized",
				},
				"group_by_source": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "For the tabular view, whether to group multi metrics into a single row by a common source. If false, each metric for each source is displayed in its own row. If true, multiple metrics for the same host will be displayed as different columns in the same row",
				},
				"invert_dynamic_legend_hover_control": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether to disable the display of the floating legend (but reenable it when the ctrl-key is pressed)",
				},
				"line_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Plot interpolation type. linear is default = ['linear', 'step-before', 'step-after', 'basis', 'cardinal', 'monotone']",
				},
				"max": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "Max value of Y-axis. Set to null or leave blank for auto",
				},
				"min": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "Min value of Y-axis. Set to null or leave blank for auto",
				},
				"num_tags": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "For the tabular view, how many point tags to display",
				},
				"plain_markdown_content": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The Markdown content for a Markdown display, in plain text. Use this field instead of markdownContent",
				},
				"show_hosts": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "For the tabular view, whether to display sources. Default: true",
				},
				"show_labels": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "For the tabular view, whether to display labels. Default: true",
				},
				"show_raw_values": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "For the tabular view, whether to display raw values. Default: false",
				},
				"sort_values_descending": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "For the tabular view, whether to display display values in descending order. Default: false",
				},
				"sparkline_decimal_precision": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "For the single stat view, the decimal precision of the displayed number ",
				},
				"sparkline_display_color": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, the color of the displayed text (when not dynamically determined). Values should be in rgba(, , ,  format ",
				},
				"sparkline_display_font_size": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, the font size of the displayed text, in percent",
				},
				"sparkline_display_horizontal_position": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, the horizontal position of the displayed text = ['MIDDLE', 'LEFT', 'RIGHT']",
				},
				"sparkline_display_postfix": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, a string to append to the displayed text",
				},
				"sparkline_display_prefix": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, a string to add before the displayed text",
				},
				"sparkline_display_value_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, whether to display the name of the query or the value of query = ['VALUE', 'LABEL']",
				},
				"sparkline_display_vertical_position": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "deprecated",
				},
				"sparkline_fill_color": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, the color of the background fill. Values should be in rgba(, , ,  format",
				},
				"sparkline_line_color": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, the color of the line. Values should be in rgba(, , ,  format",
				},

				"sparkline_size": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, a misleadingly named property. This determines whether the sparkline of the statistic is displayed in the chart BACKGROUND, BOTTOM, or NONE = ['BACKGROUND', 'BOTTOM', 'NONE']",
				},
				"sparkline_value_color_map_apply_to": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the single stat view, whether to apply dynamic color settings to the displayed TEXT or BACKGROUND = ['TEXT', 'BACKGROUND']",
				},
				"sparkline_value_color_map_colors": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "For the single stat view, a list of colors that differing query values map to. Must contain one more element than sparklineValueColorMapValuesV2. Values should be in rgba(, , ,  format",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"sparkline_value_color_map_values": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "deprecated",
					Elem:        &schema.Schema{Type: schema.TypeInt},
				},
				"sparkline_value_color_map_values_v2": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "deprecated",
					Elem:        &schema.Schema{Type: schema.TypeFloat},
				},
				"sparkline_value_text_map_text": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "For the single stat view, a list of display text values that different query values map to. Must contain one more element than sparklineValueTextMapThresholds",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"sparkline_value_text_map_thresholds": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "For the single stat view, a list of threshold boundaries for mapping different query values to display text. Must contain one less element than sparklineValueTextMapText",
					Elem:        &schema.Schema{Type: schema.TypeFloat},
				},
				"stack_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Type of stacked chart (applicable only if chart type is stacked). zero (default) means stacked from y=0. expand means Normalized from 0 to 1. wiggle means Minimize weighted changes. silhouette means to Center the Stream = ['zero', 'expand', 'wiggle', 'silhouette']",
				},
				"tag_mode": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the tabular view, which mode to use to determine which point tags to display = ['all', 'top', 'custom']",
				},
				"time_based_coloring": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Fox x-y scatterplots, whether to color more recent points as darker than older points. Default: false",
				},
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Chart Type. 'line' refers to the Line Plot, 'scatter' to the Point Plot, 'stacked-area' to the Stacked Area plot, 'table' to the Tabular View, 'scatterploy-xy' to Scatter Plot, 'markdown-widget' to the Markdown display, and 'sparkline' to the Single Stat view = ['line', 'scatterplot', 'stacked-area', 'table', 'scatterplot-xy', 'markdown-widget', 'sparkline']",
				},
				"windowing": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For the tabular view, whether to use the full time window for the query or the last X minutes = ['full', 'last']",
				},
				"window_size": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Width, in minutes, of the time window to use for last windowing ",
				},
				"xmax": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "For x-y scatterplots, max value for X-axis. Set null for auto",
				},
				"xmin": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "For x-y scatterplots, min value for X-axis. Set null for auto",
				},
				"y0_scale_si_by_1024": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Default: false. Whether to scale numerical magnitude labels for left Y-axis by 1024 in the IEC/Binary manner (instead of by 1000 like SI) ,",
				},
				"y0_unit_autoscaling": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Default: false. Whether to automatically adjust magnitude labels and units for the left Y-axis to favor smaller magnitudes and larger units",
				},
				"y1max": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "For plots with multiple Y-axes, max value for right-side Y-axis. Set null for auto",
				},
				"y1min": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "For plots with multiple Y-axes, min value for right-side Y-axis. Set null for auto",
				},
				"y1_scale_si_by_1024": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Default: false. Whether to scale numerical magnitude labels for right Y-axis by 1024 in the IEC/Binary manner (instead of by 1000 like SI)",
				},
				"y1_unit_autoscaling": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Default: false. Whether to automatically adjust magnitude labels and units for the right Y-axis to favor smaller magnitudes and larger units",
				},
				"y1_units": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "For plots with multiple Y-axes, units for right-side Y-axis ",
				},
				"ymax": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "For x-y scatterplots, max value for Y-axis. Set null for auto ",
				},
				"ymin": {
					Type:        schema.TypeFloat,
					Optional:    true,
					Description: "For x-y scatterplots, min value for Y-axis. Set null for auto",
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
				"summarization": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Summarization strategy for the chart. MEAN is default = ['MEAN', 'MEDIAN', 'MIN', 'MAX', 'SUM', 'COUNT', 'LAST', 'FIRST']",
				},
				"source":        source,
				"chart_setting": chartSetting,
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
			"display_section_table_of_contents": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"display_query_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"event_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	chart["chart_setting"] = []interface{}{buildTerraformChartSettings(wavefrontChart.ChartSettings)}
	return chart
}

func buildTerraformChartSettings(wavefrontChartSettings wavefront.ChartSetting) map[string]interface{} {
	chartSettings := map[string]interface{}{}
	chartSettings["auto_column_tags"] = wavefrontChartSettings.AutoColumnTags
	chartSettings["column_tags"] = wavefrontChartSettings.ColumnTags
	chartSettings["custom_tags"] = wavefrontChartSettings.CustomTags
	chartSettings["expected_data_spacing"] = wavefrontChartSettings.ExpectedDataSpacing
	chartSettings["fixed_legend_display_stats"] = wavefrontChartSettings.FixedLegendDisplayStats
	chartSettings["fixed_legend_enabled"] = wavefrontChartSettings.FixedLegendEnabled
	chartSettings["fixed_legend_filter_field"] = wavefrontChartSettings.FixedLegendFilterField
	chartSettings["fixed_legend_filter_limit"] = wavefrontChartSettings.FixedLegendFilterLimit
	chartSettings["fixed_legend_filter_sort"] = wavefrontChartSettings.FixedLegendFilterSort
	chartSettings["fixed_legend_hide_label"] = wavefrontChartSettings.FixedLegendHideLabel
	chartSettings["fixed_legend_position"] = wavefrontChartSettings.FixedLegendPosition
	chartSettings["fixed_legend_use_raw_stats"] = wavefrontChartSettings.FixedLegendUseRawStats
	chartSettings["group_by_source"] = wavefrontChartSettings.GroupBySource
	chartSettings["invert_dynamic_legend_hover_control"] = wavefrontChartSettings.InvertDynamicLegendHoverControl
	chartSettings["line_type"] = wavefrontChartSettings.LineType
	chartSettings["max"] = wavefrontChartSettings.Max
	chartSettings["min"] = wavefrontChartSettings.Min
	chartSettings["num_tags"] = wavefrontChartSettings.NumTags
	chartSettings["plain_markdown_content"] = wavefrontChartSettings.PlainMarkdownContent
	chartSettings["show_hosts"] = wavefrontChartSettings.ShowHosts
	chartSettings["show_labels"] = wavefrontChartSettings.ShowLabels
	chartSettings["show_raw_values"] = wavefrontChartSettings.ShowRawValues
	chartSettings["sort_values_descending"] = wavefrontChartSettings.SortValuesDescending
	chartSettings["sparkline_decimal_precision"] = wavefrontChartSettings.SparklineDecimalPrecision
	chartSettings["sparkline_display_color"] = wavefrontChartSettings.SparklineDisplayColor
	chartSettings["sparkline_display_font_size"] = wavefrontChartSettings.SparklineDisplayFontSize
	chartSettings["sparkline_display_horizontal_position"] = wavefrontChartSettings.SparklineDisplayHorizontalPosition
	chartSettings["sparkline_display_postfix"] = wavefrontChartSettings.SparklineDisplayPostfix
	chartSettings["sparkline_display_prefix"] = wavefrontChartSettings.SparklineDisplayPrefix
	chartSettings["sparkline_display_value_type"] = wavefrontChartSettings.SparklineDisplayValueType
	chartSettings["sparkline_display_vertical_position"] = wavefrontChartSettings.SparklineDisplayVerticalPosition
	chartSettings["sparkline_fill_color"] = wavefrontChartSettings.SparklineFillColor
	chartSettings["sparkline_line_color"] = wavefrontChartSettings.SparklineLineColor
	chartSettings["sparkline_size"] = wavefrontChartSettings.SparklineSize
	chartSettings["sparkline_value_color_map_apply_to"] = wavefrontChartSettings.SparklineValueColorMapApplyTo
	chartSettings["sparkline_value_color_map_colors"] = wavefrontChartSettings.SparklineValueColorMapColors
	chartSettings["sparkline_value_color_map_values"] = wavefrontChartSettings.SparklineValueColorMapValues
	chartSettings["sparkline_value_color_map_values_v2"] = wavefrontChartSettings.SparklineValueColorMapValuesV2
	chartSettings["sparkline_value_text_map_text"] = wavefrontChartSettings.SparklineValueTextMapText
	chartSettings["sparkline_value_text_map_thresholds"] = wavefrontChartSettings.SparklineValueTextMapThresholds
	chartSettings["stack_type"] = wavefrontChartSettings.StackType
	chartSettings["tag_mode"] = wavefrontChartSettings.TagMode
	chartSettings["time_based_coloring"] = wavefrontChartSettings.TimeBasedColoring
	chartSettings["type"] = wavefrontChartSettings.Type
	chartSettings["windowing"] = wavefrontChartSettings.Windowing
	chartSettings["window_size"] = wavefrontChartSettings.WindowSize
	chartSettings["xmax"] = wavefrontChartSettings.Xmax
	chartSettings["xmin"] = wavefrontChartSettings.Xmin
	chartSettings["y0_scale_si_by_1024"] = wavefrontChartSettings.Y0ScaleSIBy1024
	chartSettings["y1_scale_si_by_1024"] = wavefrontChartSettings.Y1ScaleSIBy1024
	chartSettings["y0_unit_autoscaling"] = wavefrontChartSettings.Y0UnitAutoscaling
	chartSettings["y1_unit_autoscaling"] = wavefrontChartSettings.Y1UnitAutoscaling
	chartSettings["y1_units"] = wavefrontChartSettings.Y1Units
	chartSettings["y1max"] = wavefrontChartSettings.Y1Max
	chartSettings["y1min"] = wavefrontChartSettings.Y1Min
	return chartSettings
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
		terraformChartSettings := t["chart_setting"].([]interface{})

		wavefrontCharts[i] = wavefront.Chart{
			Name:          t["name"].(string),
			Sources:       *buildSources(&terraformSources),
			Description:   t["description"].(string),
			Units:         t["units"].(string),
			Summarization: t["summarization"].(string),
			ChartSettings: *buildChartSettings(&terraformChartSettings),
		}
	}

	return &wavefrontCharts
}

// Construct a Wavefront ChartSetting
func buildChartSettings(terraformChartSettings *[]interface{}) *wavefront.ChartSetting {
	wavefrontChartSettings := &wavefront.ChartSetting{}

	t := ((*terraformChartSettings)[0]).(map[string]interface{})

	if t["auto_column_tags"] != nil {
		wavefrontChartSettings.AutoColumnTags = t["auto_column_tags"].(bool)
	}
	if t["column_tags"] != nil {
		wavefrontChartSettings.ColumnTags = t["column_tags"].(string)
	}
	if t["custom_tags"] != nil {
		for _, tag := range t["custom_tags"].([]interface{}) {
			wavefrontChartSettings.CustomTags = append(wavefrontChartSettings.CustomTags, tag.(string))
		}
	}
	if t["expected_data_spacing"] != nil {
		wavefrontChartSettings.ExpectedDataSpacing = t["expected_data_spacing"].(int)
	}
	if t["fixed_legend_display_stats"] != nil {
		for _, stat := range t["fixed_legend_display_stats"].([]interface{}) {
			wavefrontChartSettings.FixedLegendDisplayStats = append(wavefrontChartSettings.FixedLegendDisplayStats, stat.(string))
		}
	}
	if t["fixed_legend_enabled"] != nil {
		wavefrontChartSettings.FixedLegendEnabled = t["fixed_legend_enabled"].(bool)
	}
	if t["fixed_legend_filter_field"] != nil {
		wavefrontChartSettings.FixedLegendFilterField = t["fixed_legend_filter_field"].(string)
	}
	if t["fixed_legend_filter_limit"] != nil {
		wavefrontChartSettings.FixedLegendFilterLimit = t["fixed_legend_filter_limit"].(int)
	}
	if t["fixed_legend_filter_sort"] != nil {
		wavefrontChartSettings.FixedLegendFilterSort = t["fixed_legend_filter_sort"].(string)
	}
	if t["fixed_legend_hide_label"] != nil {
		wavefrontChartSettings.FixedLegendHideLabel = t["fixed_legend_hide_label"].(bool)
	}
	if t["fixed_legend_position"] != nil {
		wavefrontChartSettings.FixedLegendPosition = t["fixed_legend_position"].(string)
	}
	if t["fixed_legend_use_raw_stats"] != nil {
		wavefrontChartSettings.FixedLegendUseRawStats = t["fixed_legend_use_raw_stats"].(bool)
	}
	if t["group_by_source"] != nil {
		wavefrontChartSettings.GroupBySource = t["group_by_source"].(bool)
	}
	if t["invert_dynamic_legend_hover_control"] != nil {
		wavefrontChartSettings.InvertDynamicLegendHoverControl = t["invert_dynamic_legend_hover_control"].(bool)
	}
	if t["line_type"] != nil {
		wavefrontChartSettings.LineType = t["line_type"].(string)
	}
	if t["max"] != nil {
		wavefrontChartSettings.Max = float32(t["max"].(float64))
	}
	if t["min"] != nil {
		wavefrontChartSettings.Min = float32(t["min"].(float64))
	}
	if t["num_tags"] != nil {
		wavefrontChartSettings.NumTags = t["num_tags"].(int)
	}
	if t["plain_markdown_content"] != nil {
		wavefrontChartSettings.PlainMarkdownContent = t["plain_markdown_content"].(string)
	}
	if t["show_hosts"] != nil {
		wavefrontChartSettings.ShowHosts = t["show_hosts"].(bool)
	}
	if t["show_labels"] != nil {
		wavefrontChartSettings.ShowLabels = t["show_labels"].(bool)
	}
	if t["show_raw_values"] != nil {
		wavefrontChartSettings.ShowRawValues = t["show_raw_values"].(bool)
	}
	if t["sort_values_descending"] != nil {
		wavefrontChartSettings.SortValuesDescending = t["sort_values_descending"].(bool)
	}
	if t["sparkline_decimal_precision"] != nil {
		wavefrontChartSettings.SparklineDecimalPrecision = t["sparkline_decimal_precision"].(int)
	}
	if t["sparkline_display_color"] != nil {
		wavefrontChartSettings.SparklineDisplayColor = t["sparkline_display_color"].(string)
	}
	if t["sparkline_display_font_size"] != nil {
		wavefrontChartSettings.SparklineDisplayFontSize = t["sparkline_display_font_size"].(string)
	}
	if t["sparkline_display_horizontal_position"] != nil {
		wavefrontChartSettings.SparklineDisplayHorizontalPosition = t["sparkline_display_horizontal_position"].(string)
	}
	if t["sparkline_display_postfix"] != nil {
		wavefrontChartSettings.SparklineDisplayPostfix = t["sparkline_display_postfix"].(string)
	}
	if t["sparkline_display_prefix"] != nil {
		wavefrontChartSettings.SparklineDisplayPrefix = t["sparkline_display_prefix"].(string)
	}
	if t["sparkline_display_value_type"] != nil {
		wavefrontChartSettings.SparklineDisplayValueType = t["sparkline_display_value_type"].(string)
	}
	if t["sparkline_display_vertical_position"] != nil {
		wavefrontChartSettings.SparklineDisplayVerticalPosition = t["sparkline_display_vertical_position"].(string)
	}
	if t["sparkline_fill_color"] != nil {
		wavefrontChartSettings.SparklineFillColor = t["sparkline_fill_color"].(string)
	}
	if t["sparkline_line_color"] != nil {
		wavefrontChartSettings.SparklineLineColor = t["sparkline_line_color"].(string)
	}
	if t["sparkline_size"] != nil {
		wavefrontChartSettings.SparklineSize = t["sparkline_size"].(string)
	}
	if t["sparkline_value_color_map_apply_to"] != nil {
		wavefrontChartSettings.SparklineValueColorMapApplyTo = t["sparkline_value_color_map_apply_to"].(string)
	}
	if t["sparkline_value_color_map_colors"] != nil {
		for _, v := range t["sparkline_value_color_map_colors"].([]interface{}) {
			wavefrontChartSettings.SparklineValueColorMapColors = append(wavefrontChartSettings.SparklineValueColorMapColors, v.(string))
		}
	}
	if t["sparkline_value_color_map_values"] != nil {
		for _, v := range t["sparkline_value_color_map_values"].([]interface{}) {

			wavefrontChartSettings.SparklineValueColorMapValues = append(wavefrontChartSettings.SparklineValueColorMapValues, v.(int))
		}
	}
	if t["sparkline_value_text_map_text"] != nil {
		for _, v := range t["sparkline_value_text_map_text"].([]interface{}) {
			wavefrontChartSettings.SparklineValueTextMapText = append(wavefrontChartSettings.SparklineValueTextMapText, v.(string))
		}
	}
	if t["sparkline_value_text_map_thresholds"] != nil {
		for _, v := range t["sparkline_value_text_map_thresholds"].([]interface{}) {
			wavefrontChartSettings.SparklineValueTextMapThresholds = append(wavefrontChartSettings.SparklineValueTextMapThresholds, float32(v.(float64)))
		}
	}
	if t["sparkline_value_color_map_values_v2"] != nil {
		for _, v := range t["sparkline_value_color_map_values_v2"].([]interface{}) {
			wavefrontChartSettings.SparklineValueColorMapValuesV2 = append(wavefrontChartSettings.SparklineValueColorMapValuesV2, float32(v.(float64)))
		}
	}

	if t["stack_type"] != nil {
		wavefrontChartSettings.StackType = t["stack_type"].(string)
	}
	if t["tag_mode"] != nil {
		wavefrontChartSettings.TagMode = t["tag_mode"].(string)
	}
	if t["time_based_coloring"] != nil {
		wavefrontChartSettings.TimeBasedColoring = t["time_based_coloring"].(bool)
	}
	if t["type"] != nil {
		wavefrontChartSettings.Type = t["type"].(string)
	}
	if t["windowing"] != nil {
		wavefrontChartSettings.Windowing = t["windowing"].(string)
	}
	if t["window_size"] != nil {
		wavefrontChartSettings.WindowSize = t["window_size"].(int)
	}
	if t["xmax"] != nil {
		wavefrontChartSettings.Xmax = float32(t["xmax"].(float64))
	}
	if t["xmin"] != nil {
		wavefrontChartSettings.Xmin = float32(t["xmin"].(float64))
	}
	if t["y0_scale_si_by_1024"] != nil {
		wavefrontChartSettings.Y0ScaleSIBy1024 = t["y0_scale_si_by_1024"].(bool)
	}
	if t["y0_unit_autoscaling"] != nil {
		wavefrontChartSettings.Y0UnitAutoscaling = t["y0_unit_autoscaling"].(bool)
	}
	if t["y1max"] != nil {
		wavefrontChartSettings.Y1Max = float32(t["y1max"].(float64))
	}
	if t["y1min"] != nil {
		wavefrontChartSettings.Y1Min = float32(t["y1min"].(float64))
	}
	if t["y1_scale_si_by_1024"] != nil {
		wavefrontChartSettings.Y1ScaleSIBy1024 = t["y1_scale_si_by_1024"].(bool)
	}
	if t["y1_unit_autoscaling"] != nil {
		wavefrontChartSettings.Y1UnitAutoscaling = t["y1_unit_autoscaling"].(bool)
	}
	if t["y1_units"] != nil {
		wavefrontChartSettings.Y1Units = t["y1_units"].(string)
	}
	if t["ymax"] != nil {
		wavefrontChartSettings.Ymax = float32(t["ymax"].(float64))
	}
	if t["ymin"] != nil {
		wavefrontChartSettings.Ymin = float32(t["ymin"].(float64))
	}

	return wavefrontChartSettings
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
	eventFilterType := "BYCHART"
	if e, ok := d.GetOk("event_filter_type"); ok {
		eventFilterType = e.(string)
	}

	displayTOC := false
	if toc, ok := d.GetOk("display_section_table_of_contents"); ok {
		displayTOC = toc.(bool)
	}

	displayQP := false
	if qp, ok := d.GetOk("display_query_parameters"); ok {
		displayQP = qp.(bool)
	}

	return &wavefront.Dashboard{
		Name:                          d.Get("name").(string),
		ID:                            d.Get("url").(string),
		Tags:                          tags,
		Description:                   d.Get("description").(string),
		Url:                           d.Get("url").(string),
		Sections:                      *buildSections(&terraformSections),
		ParameterDetails:              *buildParameterDetails(&terraformParams),
		EventFilterType:               eventFilterType,
		DisplaySectionTableOfContents: displayTOC,
		DisplayQueryParameters:        displayQP,
	}, nil
}

// Create a Terraform Dashboard
func resourceDashboardCreate(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()
	dashboard, err := buildDashboard(d)

	if err != nil {
		return fmt.Errorf("failed to parse dashboard, %s", err)
	}

	err = dashboards.Create(dashboard)
	if err != nil {
		return fmt.Errorf("failed to create dashboard, %s", err)
	}
	d.SetId(dashboard.ID)

	return resourceDashboardRead(d, m)
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

	d.Set("display_section_table_of_contents", dash.DisplaySectionTableOfContents)
	d.Set("display_query_parameters", dash.DisplayQueryParameters)

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
	return resourceDashboardRead(d, m)
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
