package configs

import (
	"fmt"
	"log/slog"
	"path"
	"time"

	"gopkg.in/ini.v1"
)

const (
	configDir       = "config"
	configFilename  = "config.ini"
	dateLayout      = "060102"
	sectionFolders  = "folders"
	inputDirKey     = "inputDir"
	outputDirKey    = "outputDir"
	sectionDate     = "date"
	dateKey         = "date"
	sectionApp      = "app"
	logLevelKey     = "logLevel"
	logLevelFileKey = "logLevelFile"
	logsDirKey      = "logsDir"
	dbFile          = "dbFile"
)

type Config struct {
	InputDir     string
	OutputDir    string
	LogsDir      string
	Date         string
	LogLevel     slog.Level
	LogLevelFile slog.Level
	DBFile       string
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

	logLevelRaw := iniCfg.Section(sectionApp).Key(logLevelKey)
	if logLevelRaw == nil {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", logLevelKey, sectionApp)
	}

	logLevelFileRaw := iniCfg.Section(sectionApp).Key(logLevelFileKey)
	if logLevelFileRaw == nil {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", logLevelFileKey, sectionApp)
	}

	logsDir := iniCfg.Section(sectionApp).Key(logsDirKey)
	if logsDir == nil {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", logsDirKey, sectionApp)
	}

	dbFileRaw := iniCfg.Section(sectionApp).Key(dbFile)
	if dbFileRaw == nil {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", dbFile, sectionApp)
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

	if logLevelRaw.String() == "" {
		return nil, fmt.Errorf("log level is empty")
	}

	logLevel := slog.LevelInfo
	if logLevelRaw.String() != "" {
		e := logLevel.UnmarshalText([]byte(logLevelRaw.String()))
		if e != nil {
			return nil, fmt.Errorf("invalid log level: %v", err)
		}
	}

	logLevelFile := slog.LevelWarn
	if logLevelFileRaw.String() != "" {
		e := logLevel.UnmarshalText([]byte(logLevelRaw.String()))
		if e != nil {
			return nil, fmt.Errorf("invalid log level: %v", err)
		}
	}

	if logsDir.String() == "" {
		return nil, fmt.Errorf("logs dir is empty")
	}

	if dbFileRaw.String() == "" {
		return nil, fmt.Errorf("logs dir is empty")
	}

	return &Config{
		InputDir:     inputDir.String(),
		OutputDir:    outputDir.String(),
		Date:         date.Format(dateLayout),
		LogsDir:      logsDir.String(),
		LogLevel:     logLevel,
		LogLevelFile: logLevelFile,
		DBFile:       dbFileRaw.String(),
	}, nil
}
