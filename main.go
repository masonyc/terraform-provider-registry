package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-registry/registry"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: "serko.com/serko/registry",
		ProviderFunc: registry.Provider,
	}

	plugin.Serve(opts)
}
