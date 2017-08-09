package wavefront

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type wavefrontConfig struct {
	address string
	token   string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAVEFRONT_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wavefront_alert": resourceAgent(),
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
