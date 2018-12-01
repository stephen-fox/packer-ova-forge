package ovaforge

import (
	"os"
	"strings"
)

type forgeArtifacts struct {
	filePaths []string
}

func (*forgeArtifacts) BuilderId() string {
	return "stuff"
}

func (o *forgeArtifacts) Files() []string {
	return o.filePaths
}

func (o *forgeArtifacts) Id() string {
	return "ova"
}

func (o *forgeArtifacts) String() string {
	return "OVA related files: " + strings.Join(o.filePaths, ", ")
}

func (o *forgeArtifacts) State(name string) interface{} {
	return nil
}

func (o *forgeArtifacts) Destroy() error {
	for _, filePath := range o.filePaths {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
