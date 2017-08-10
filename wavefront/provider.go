package wavefront_plugin

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

type wavefrontConfig struct {
	address string
	token   string
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAVEFRONT_ADDRESS", ""),
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAVEFRONT_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wavefront_alert": resourceAlert(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return &wavefrontConfig{
		address: d.Get("address").(string),
		token:   d.Get("token").(string),
	}, nil

}
