package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	wavefront_provider "github.com/spaceapegames/terraform-provider-wavefront/wavefront"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return wavefront_provider.Provider()
		},
	})
}
