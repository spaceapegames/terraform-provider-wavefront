package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/spaceapegames/terraform-provider-wavefront/wavefront"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: wavefront_plugin.Provider})
}
