package configs

import (
	"fmt"
	"path"
	"time"

	"gopkg.in/ini.v1"
)

const (
	configDir      = "config"
	configFilename = "config.ini"
	dateLayout     = "060102"
	sectionFolders = "folders"
	inputDirKey    = "inputDir"
	outputDirKey   = "outputDir"
	sectionDate    = "date"
	dateKey        = "date"
)

type Config struct {
	InputDir  string
	OutputDir string
	Date      string
}

func MustLoad() (*Config, error) {
	iniCfg, err := ini.Load(path.Join(configDir, configFilename))
	if err != nil {
		return nil, err
	}

	inputDir := iniCfg.Section(sectionFolders).Key(inputDirKey)
	if inputDir == nil {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", inputDirKey, sectionFolders)
	}

	outputDir := iniCfg.Section(sectionFolders).Key(outputDirKey)
	if outputDir == nil {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", outputDirKey, sectionFolders)
	}

	dateRaw := iniCfg.Section(sectionDate).Key(dateKey)
	if dateRaw == nil {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", dateKey, sectionDate)
	}

	date, err := time.Parse(dateLayout, dateRaw.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %v", err)
	}

	if inputDir.String() == "" {
		return nil, fmt.Errorf("input dir path is empty")
	}

	if outputDir.String() == "" {
		return nil, fmt.Errorf("output dir path is empty")
	}

	if dateRaw.String() == "" {
		return nil, fmt.Errorf("date is empty")
	}

	return &Config{
		InputDir:  inputDir.String(),
		OutputDir: outputDir.String(),
		Date:      date.Format(dateLayout),
	}, nil
}
