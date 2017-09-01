package wavefront_plugin

import (
	"fmt"
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
				"string_key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"string_value": {
					Type:     schema.TypeString,
					Required: true,
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
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// Construct a Terraform ParameterDetail
func buildTerraformParameterDetail(wavefrontParamDetail wavefront.ParameterDetail, name string) map[string]interface{} {
	parameterDetail := map[string]interface{}{}

	for k, v := range wavefrontParamDetail.ValuesToReadableStrings {
		parameterDetail["string_key"] = k
		parameterDetail["string_value"] = v
	}
	parameterDetail["name"] = name
	parameterDetail["label"] = wavefrontParamDetail.Label
	parameterDetail["parameter_type"] = wavefrontParamDetail.ParameterType
	parameterDetail["hide_from_view"] = wavefrontParamDetail.HideFromView
	parameterDetail["default_value"] = wavefrontParamDetail.DefaultValue
	return parameterDetail
}

// Construct a Terraform Section
func buildTerraformSection(wavefront_section wavefront.Section) map[string]interface{} {
	section := map[string]interface{}{}
	section["name"] = wavefront_section.Name
	rows := []map[string]interface{}{}
	for _, wavefront_row := range wavefront_section.Rows {
		rows = append(rows, buildTerraformRow(wavefront_row))
	}
	section["row"] = rows

	return section
}

// Construct a Wavefront Row
func buildTerraformRow(wavefront_row wavefront.Row) map[string]interface{} {
	row := map[string]interface{}{}

	charts := []map[string]interface{}{}
	for _, wavefront_row := range wavefront_row.Charts {
		charts = append(charts, buildTerraformChart(wavefront_row))
	}
	row["chart"] = charts

	return row
}

// Construct a Wavefront Chart
func buildTerraformChart(wavefront_chart wavefront.Chart) map[string]interface{} {
	chart := map[string]interface{}{}
	chart["name"] = wavefront_chart.Name
	chart["description"] = wavefront_chart.Description

	chart["units"] = wavefront_chart.Units
	sources := []map[string]interface{}{}
	for _, wavefront_source := range wavefront_chart.Sources {
		sources = append(sources, buildTerraformSource(wavefront_source))
	}
	chart["source"] = sources

	return chart
}

// Construct a Wavefront Source
func buildTerraformSource(wavefront_source wavefront.Source) map[string]interface{} {
	source := map[string]interface{}{}
	source["name"] = wavefront_source.Name
	source["query"] = wavefront_source.Query

	return source
}

// Construct a Wavefront Section
func buildSections(terraformSections *[]interface{}) *[]wavefront.Section {
	wavefrontSections := make([]wavefront.Section, len(*terraformSections))

	for i, t_ := range *terraformSections {
		t := t_.(map[string]interface{})

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

	for i, t_ := range *terraformRows {
		t := t_.(map[string]interface{})

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

	for i, t_ := range *terraformCharts {
		t := t_.(map[string]interface{})

		terraformSources := t["source"].([]interface{})

		wavefrontCharts[i] = wavefront.Chart{
			Name:        t["name"].(string),
			Sources:     *buildSources(&terraformSources),
			Description: t["description"].(string),
			Units:       t["units"].(string),
		}
	}

	return &wavefrontCharts
}

// Construct a Wavefront Source
func buildSources(terrafromSources *[]interface{}) *[]wavefront.Source {
	wavefrontSources := make([]wavefront.Source, len(*terrafromSources))

	for i, t_ := range *terrafromSources {
		t := t_.(map[string]interface{})

		wavefrontSources[i] = wavefront.Source{
			Name:  t["name"].(string),
			Query: t["query"].(string),
		}
	}

	return &wavefrontSources
}

// Construct a Wavefront ParameterDetail
func buildParameterDetails(terraformParams *[]interface{}) *map[string]wavefront.ParameterDetail {
	wavefrontParams := map[string]wavefront.ParameterDetail{}

	for _, t_ := range *terraformParams {
		t := t_.(map[string]interface{})

		wavefrontParams[t["name"].(string)] = wavefront.ParameterDetail{
			Label:         t["label"].(string),
			DefaultValue:  t["default_value"].(string),
			HideFromView:  t["hide_from_view"].(bool),
			ParameterType: t["parameter_type"].(string),
			ValuesToReadableStrings: map[string]string{
				t["string_key"].(string): t["string_value"].(string),
			},
		}
	}

	return &wavefrontParams
}

// Construct a Wavefront Dashboard
func buildDashboard(d *schema.ResourceData) (*wavefront.Dashboard, error) {

	var tags []string
	for _, tag := range d.Get("tags").([]interface{}) {
		tags = append(tags, tag.(string))
	}

	terraformSections := d.Get("section").([]interface{})
	terraformParams := d.Get("parameter_details").([]interface{})

	return &wavefront.Dashboard{
		Name:             d.Get("name").(string),
		ID:               d.Get("url").(string),
		Tags:             tags,
		Description:      d.Get("description").(string),
		Url:              d.Get("url").(string),
		Sections:         *buildSections(&terraformSections),
		ParameterDetails: *buildParameterDetails(&terraformParams),
	}, nil
}

// Create a Terraform Dashboard
func resourceDashboardCreate(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()
	dashboard, err := buildDashboard(d)
	if err != nil {
		fmt.Errorf("Failed to parse dashboard, %s", err)
	}

	err = dashboards.Create(dashboard)
	if err != nil {
		fmt.Errorf("Failed to create dashboard, %s", err)
	}
	d.SetId(dashboard.ID)

	return nil
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
		// dashboard no longer exists
		d.SetId("")
	}

	// Use the Wavefront url as the Terraform ID
	d.SetId(dash.ID)
	d.Set("name", dash.Name)
	d.Set("description", dash.Description)
	d.Set("url", dash.Url)

	sections := []map[string]interface{}{}
	for _, wavefrontSection := range dash.Sections {
		sections = append(sections, buildTerraformSection(wavefrontSection))
	}
	d.Set("section", sections)

	parameterDetails := []map[string]interface{}{}

	for k, v := range dash.ParameterDetails {
		parameterDetails = append(parameterDetails, buildTerraformParameterDetail(v, k))
	}

	d.Set("parameter_details", parameterDetails)
	d.Set("tags", dash.Tags)

	return nil
}

func resourceDashboardUpdate(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()

	a, err := buildDashboard(d)
	if err != nil {
		fmt.Errorf("Failed to parse dashboard, %s", err)
	}

	// Update the dashboard on Wavefront
	err = dashboards.Update(a)
	if err != nil {
		return fmt.Errorf("Error Updating Dashboard %s. %s", d.Get("name"), err)
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
		return fmt.Errorf("Error finding Wavefront Dashboard %s. %s", d.Id(), err)
	}

	// Delete the Dashbaord
	err = dashboards.Delete(&dash)
	if err != nil {
		return fmt.Errorf("Failed to delete Dashboard %s. %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
