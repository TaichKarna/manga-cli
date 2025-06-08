package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const configFileName = "config.json"
const configDirName = ".manga-cli"

type Config map[string]any

var ValidConfigOptions = map[string]struct {
	Description string
	Default     any
}{
	"path": {Description: "Path where downloaded manga is stored", Default: "~/Pictures/manga-cli"},
	"viewer":        {Description: "External image viewer (e.g., viu, feh, imv, sxiv)", Default: "viu"},
	"language":      {Description: "Preferred language for manga", Default: "en"},
	"width": {
    	Description: "Default image width for terminal viewer",
    	Default:     60,
	},
	"height": {
    	Description: "Default image height for terminal viewer",
    	Default:     40,
	},

}


func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}

func toRelativeHomePath(path string) string {
	usr, err := user.Current()
	if err != nil {
		return path 
	}
	home := usr.HomeDir

	path = expandPath(path) 

	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	rel, err := filepath.Rel(home, absPath)
	if err != nil {
		return path
	}

	if !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(filepath.Join("~", rel))
	}
	return path 
}


func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine user home directory")
	}
	return filepath.Join(home, configDirName, configFileName)
}

func ensureConfigDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(home, configDirName), 0755)
}

func LoadConfig() (Config, error) {
	path := getConfigFilePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaults := Config{}
		for key, meta := range ValidConfigOptions {
			defaults[key] = meta.Default
		}
		if err := saveConfig(defaults); err != nil {
			return nil, err
		}
		return defaults, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}


func saveConfig(cfg Config) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}
	path := getConfigFilePath()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SetConfigOption(key string, value any) error {
	if _, ok := ValidConfigOptions[key]; !ok {
		return fmt.Errorf("invalid config key '%s'. Use `manga-cli config list` to see valid keys", key)
	}

	if key == "path" {
		if strVal, ok := value.(string); ok {
			value = toRelativeHomePath(strVal)
		}
	}

	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	cfg[key] = value
	return saveConfig(cfg)
}


func GetConfigOption(key string) (any, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	val, ok := cfg[key]
	if !ok {
		def, defOk := ValidConfigOptions[key]
		if !defOk {
			return nil, fmt.Errorf("config key '%s' not found and no default available", key)
		}
		val = def.Default
	}

	if key == "path" {
		if strVal, ok := val.(string); ok {
			val = expandPath(strVal)
		}
	}

	return val, nil
}

func GetAllConfig() (Config, error) {
	return LoadConfig()
}



