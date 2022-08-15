package main

import (
	"PoC.RegistryState/registry"
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := providerserver.ServeOpts{
		Address: "serko.com/serko/registry",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), registry.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
