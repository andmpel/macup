package macup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Terminal color codes and constants
const (
	k_timeout    = 5 * time.Second          // Timeout for HTTP requests
	k_testURL    = "https://www.google.com" // URL to test internet connection
	k_gemCmdPath = "/usr/bin/gem"           // Path to the gem command
	k_configFile = ".macup.json"            // Configuration file name
)

// checkCommand checks if a command exists in `PATH`, print warning if not.
func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}
	return true
}

// runCommand runs a shell command and directs its output to .
func runCommand(name string, args ...string) (string, error) {
	// Allow only specific commands
	allowedCommands := map[string]bool{
		"brew":           true,
		"code":           true,
		"gem":            true,
		"npm":            true,
		"yarn":           true,
		"cargo":          true,
		"mas":            true,
		"softwareupdate": true,
	}

	if !allowedCommands[name] {
		return "", fmt.Errorf("command not allowed: %s", name)
	}
	// Optionally validate arguments (e.g., no special characters)
	for _, arg := range args {
		if strings.ContainsAny(arg, "&|;$><") {
			return "", fmt.Errorf("invalid argument: %s", arg)
		}
	}

	cmd := exec.Command(name, args...)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

// CheckInternet checks for internet connectivity by making an HTTP request.
func CheckInternet() bool {
	client := http.Client{
		Timeout: k_timeout,
	}

	resp, err := client.Get(k_testURL)
	if err != nil {
		fmt.Printf("⚠️ No Internet Connection!!!")
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// homeDir returns the user's home directory.
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return home
}

// configPath returns the full path to the configuration file.
func configPath() string {
	return filepath.Join(homeDir(), k_configFile)
}

// LoadConfig loads the user's selections from the configuration file.
func LoadConfig() (*Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SaveConfig saves the user's selections to the configuration file.
func (c *Config) SaveConfig() error {
	path := configPath()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
