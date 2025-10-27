package palette

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	CONFIG_DIR  = ".config/densendither"
	CONFIG_FILE = "conf.json"
)

// getConfigPath returns the full path to the config file with expanded home directory
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, CONFIG_DIR, CONFIG_FILE), nil
}

type Config map[string][]string

// takes in a string color as input, returns true if the color is a valid hex color, returns false otherwise
func validateColor(color string) bool {
	var hexPattern = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
	return hexPattern.MatchString(color)
}

// readConfig reads the JSON configuration file
func readConfig() (Config, error) {
	filename, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(Config), nil
		}
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// writeConfig writes the configuration back to the JSON file
func writeConfig(config Config) error {
	filename, err := getConfigPath()
	if err != nil {
		return err
	}
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// takes n colors as input, and stores them as a Palette struct, this struct is returned and the palette is loaded into the config file as part of the json
func AddPalette(name string, colors ...string) (Palette, error) {
	// validate input
	if name == "" {
		return Palette{}, errors.New("palette name cannot be empty")
	}
	if len(colors) == 0 {
		return Palette{}, errors.New("at least one color is required")
	}
	if len(colors) > 10 {
		return Palette{}, errors.New("maximum 10 colors allowed per palette")
	}

	// create the palette
	palette_colors := make([]string, 0, len(colors))
	for i, color := range colors {
		if !validateColor(color) {
			return Palette{}, fmt.Errorf("invalid color at position %d: %s", i, color)
		}
		palette_colors = append(palette_colors, color)
	}

	// add the palette to the file
	config, err := readConfig()
	if err != nil {
		return Palette{}, err
	}
	config[name] = palette_colors
	err = writeConfig(config)
	if err != nil {
		return Palette{}, err
	}

	// return palette
	return Palette{name, palette_colors}, nil
}

func RemovePalette(name string) error {
	// validate input
	if name == "" {
		return errors.New("palette name cannot be empty")
	}

	// loading palette
	config, err := readConfig()
	if err != nil {
		return err
	}

	// check if palette exists
	if _, exists := config[name]; !exists {
		return fmt.Errorf("palette '%s' not found", name)
	}

	// removing palette from map
	delete(config, name)

	// rewriting map
	err = writeConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func AddToPalette(name string, colors ...string) (Palette, error) {
	// validate input
	if name == "" {
		return Palette{}, errors.New("palette name cannot be empty")
	}
	if len(colors) == 0 {
		return Palette{}, errors.New("at least one color is required")
	}

	// validate all colors first
	for i, color := range colors {
		if !validateColor(color) {
			return Palette{}, fmt.Errorf("invalid color at position %d: %s", i, color)
		}
	}

	// load config
	config, err := readConfig()
	if err != nil {
		return Palette{}, err
	}

	// check if palette exists
	existingColors, exists := config[name]
	if !exists {
		return Palette{}, fmt.Errorf("palette '%s' not found", name)
	}

	// check if adding colors would exceed limit
	if len(existingColors)+len(colors) > 10 {
		return Palette{}, fmt.Errorf("adding %d colors would exceed maximum of 10 colors per palette", len(colors))
	}

	// append new colors to existing palette
	existingColors = append(existingColors, colors...)
	config[name] = existingColors

	// write config back
	err = writeConfig(config)
	if err != nil {
		return Palette{}, err
	}

	// return updated palette
	return Palette{name, existingColors}, nil
}

// GetPalette loads a palette from the config file by name
func GetPalette(name string) (Palette, error) {
	// validate input
	if name == "" {
		return Palette{}, errors.New("palette name cannot be empty")
	}

	// load config
	config, err := readConfig()
	if err != nil {
		return Palette{}, err
	}

	// check if palette exists
	colors, exists := config[name]
	if !exists {
		return Palette{}, fmt.Errorf("palette '%s' not found", name)
	}

	// return palette
	return Palette{name, colors}, nil
}

// ListPalettes returns a map of all palette names and their colors from the config file
func ListPalettes() (map[string][]string, error) {
	// load config
	config, err := readConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
