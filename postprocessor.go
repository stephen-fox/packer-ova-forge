package ovaforge

import (
	"errors"
	"path"
	"strings"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/mapstructure"
	"github.com/stephen-fox/ovaify"
	"github.com/stephen-fox/vmwareify"
)

type PostProcessor struct {
	config Configuration
}

func (o *PostProcessor) Configure(i ...interface{}) error {
	var config Configuration

	decodeErr := mapstructure.Decode(i, &config)
	if decodeErr != nil && !strings.HasSuffix(decodeErr.Error(), "expected a map, got 'slice'") {
		return decodeErr
	}

	o.config = config

	return nil
}

func (o *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	result := &forgeArtifacts{}
	var common []string
	var ovfFilePaths []string

	for _, filePath := range artifact.Files() {
		if strings.HasSuffix(filePath, ".ovf") {
			outputOvfFilePath, err := vmwareifyOvf(filePath, ui)
			if err != nil {
				return &forgeArtifacts{}, false, err
			}

			ovfFilePaths = append(ovfFilePaths, filePath)
			ovfFilePaths = append(ovfFilePaths, outputOvfFilePath)
			result.filePaths = append(result.filePaths, outputOvfFilePath)
		} else {
			common = append(common, filePath)
		}
	}

	if len(ovfFilePaths) == 0 {
		return &forgeArtifacts{}, false, errors.New("No .ovf artifacts were provided")
	}

	for _, ovfFilePath := range ovfFilePaths {
		ovaFilePath, err := createOva(ovfFilePath, common, ui)
		if err != nil {
			return &forgeArtifacts{}, false, err
		}

		result.filePaths = append(result.filePaths, ovaFilePath)
	}

	return result, true, nil
}

type Configuration struct {
	Temp bool `mapstructure:"temp"`
}

func (o *Configuration) Validate() error {
	return nil
}

func vmwareifyOvf(filePath string, ui packer.Ui) (string, error) {
	ui.Message("VMWareifying '" + filePath + "'...")

	inputFilename := path.Base(filePath)
	outputFilePath := path.Join(path.Dir(filePath), pathWithoutExtension(inputFilename) + "-vmware.ovf")

	err := vmwareify.BasicConvert(filePath, outputFilePath)
	if err != nil {
		return "", err
	}

	ui.Message("Finished VMWareifying .ovf at '" + outputFilePath + "'")

	return outputFilePath, nil
}

func createOva(ovfFilePath string, files []string, ui packer.Ui) (string, error) {
	outputPath := pathWithoutExtension(ovfFilePath) + ".ova"

	ui.Message("Creating .ova for '" + ovfFilePath + "' with files '" +
		strings.Join(files, ", ") + "'...")

	config := ovaify.OvaConfig{
		OvfFilePath:        ovfFilePath,
		FilePathsToInclude: files,
		OutputFilePath:     outputPath,
	}

	err := ovaify.ConvertOvfToOva(config)
	if err != nil {
		return "", err
	}

	ui.Message("Finished creating .ova at '" + outputPath + "'")

	return outputPath, nil
}

func pathWithoutExtension(filename string) string {
	index := strings.LastIndex(filename, ".")

	if index > 0 {
		return filename[:index]
	}

	return ""
}
