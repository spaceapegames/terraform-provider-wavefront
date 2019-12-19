package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func resourceDashboardJson() *schema.Resource {
	return &schema.Resource{
		Create: resourceDashboardJsonCreate,
		Read:   resourceDashboardJsonRead,
		Update: resourceDashboardJsonUpdate,
		Delete: resourceDashboardJsonDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"dashboard_json": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateDashboardJson,
				StateFunc:    NormalizeDashboardJson,
			},
		},
	}

}

func buildDashboardJson(d *schema.ResourceData) (*wavefront.Dashboard, error) {
	var dashboard wavefront.Dashboard
	dashboardJsonString := d.Get("dashboard_json").(string)
	// json is already validated during resource Validation
	_ = dashboard.UnmarshalJSON([]byte(dashboardJsonString))

	// set url name as the resource ID
	dashboard.ID = dashboard.Url
	return &dashboard, nil
}

func resourceDashboardJsonRead(d *schema.ResourceData, m interface{}) error {
	dashboards := m.(*wavefrontClient).client.Dashboards()
	dash := wavefront.Dashboard{
		ID: d.Id(),
	}
	err := dashboards.Get(&dash)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Wavefront Dashboard %s. %s", d.Id(), err)
		}
	}
	bytes, err := dash.MarshalJSON()
	// Use the Wavefront url as the Terraform ID
	d.SetId(dash.ID)
	err = d.Set("dashboard_json", NormalizeDashboardJson(string(bytes)))
	if err != nil {
		return fmt.Errorf("failed to set dashboard json %s. %s", d.Id(), err)
	}
	return nil
}

func resourceDashboardJsonCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Create Wavefront Dashboard %s", d.Id())
	dashboards := m.(*wavefrontClient).client.Dashboards()
	dashboard, err := buildDashboardJson(d)

	if err != nil {
		return fmt.Errorf("failed to parse dashboard, %s", err)
	}

	err = dashboards.Create(dashboard)
	if err != nil {
		return fmt.Errorf("failed to create dashboard, %s", err)
	}
	d.SetId(dashboard.ID)
	log.Printf("[INFO] Wavefront Dashboard %s Created", d.Id())
	return resourceDashboardJsonRead(d, m)
}

func resourceDashboardJsonUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Wavefront Dashboard %s", d.Id())
	dashboards := m.(*wavefrontClient).client.Dashboards()
	dashboard, err := buildDashboardJson(d)

	if err != nil {
		return fmt.Errorf("failed to parse dashboard, %s", err)
	}

	err = dashboards.Update(dashboard)
	if err != nil {
		return fmt.Errorf("failed to create dashboard, %s", err)
	}

	log.Printf("[INFO] Wavefront Dashboard %s Updated", d.Id())
	return resourceDashboardJsonRead(d, m)
}

func resourceDashboardJsonDelete(d *schema.ResourceData, m interface{}) error {
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

func ValidateDashboardJson(val interface{}, key string) ([]string, []error) {
	dashboardJsonString := val.(string)
	var dashboard wavefront.Dashboard
	err := dashboard.UnmarshalJSON([]byte(dashboardJsonString))
	if err != nil {
		return nil, []error{err}
	}
	return nil, nil
}

func NormalizeDashboardJson(val interface{}) string {
	dashboardJsonString := val.(string)
	var dashboard wavefront.Dashboard
	_ = dashboard.UnmarshalJSON([]byte(dashboardJsonString))

	// set url name as the resource ID
	dashboard.ID = dashboard.Url

	// remove keys which are not needed for diff
	dashboard.CreatedEpochMillis = 0
	dashboard.UpdatedEpochMillis = 0
	dashboard.CreatorId = ""
	dashboard.UpdaterId = ""
	dashboard.Customer = ""
	dashboard.ViewsLastDay = 0
	dashboard.ViewsLastWeek = 0
	dashboard.ViewsLastMonth = 0
	dashboard.NumCharts = 0
	dashboard.NumFavorites = 0
	dashboard.Favorite = false

	ret, _ := dashboard.MarshalJSON()
	return string(ret)
}
