package wavefront_plugin

import (
	"fmt"

	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type wavefrontClient struct {
	client wavefront.Client
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAVEFRONT_ADDRESS", ""),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAVEFRONT_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wavefront_alert":          resourceAlert(),
			"wavefront_dashboard":      resourceDashboard(),
			"wavefront_dashboard_json": resourceDashboardJson(),
			"wavefront_derived_metric": resourceDerivedMetric(),
			"wavefront_alert_target":   resourceTarget(),
			"wavefront_user":           resourceUser(),
			"wavefront_user_group":     resourceUserGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &wavefront.Config{
		Address: d.Get("address").(string),
		Token:   d.Get("token").(string),
	}
	wFClient, err := wavefront.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to configure Wavefront Client %s", err)
	}
	return &wavefrontClient{
		client: *wFClient,
	}, nil

}
