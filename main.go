package main

import (
	"github.com/WavefrontHQ/terraform-provider-wavefront/wavefront"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return wavefront_plugin.Provider()
		},
	})
}
