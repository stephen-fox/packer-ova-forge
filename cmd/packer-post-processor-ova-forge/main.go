package main

import (
	"log"

	"github.com/hashicorp/packer/packer/plugin"
	"github.com/stephen-fox/packer-post-processor-ova-forge"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		log.Fatal(err.Error())
	}

	server.RegisterPostProcessor(&ovaforge.PostProcessor{})

	server.Serve()
}
