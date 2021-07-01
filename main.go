package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/ebu/terraform-provider-mcma/mcma"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mcma.Provider,
	})
}
