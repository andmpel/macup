package macup

import (
	"os/exec"
	"strings"
)

// UpdateBrew updates Homebrew formulas and perform diagnostics.
func UpdateBrew() (string, error) {
	var output string
	output += "Updating Brew Formulas\n"
	if checkCommand("brew") {
		out, err := runCommand("brew", "update")
		output += out
		if err != nil {
			return output, err
		}
		out, err = runCommand("brew", "upgrade")
		output += out
		if err != nil {
			return output, err
		}
		out, err = runCommand("brew", "cleanup", "-s")
		output += out
		if err != nil {
			return output, err
		}
		output += "Brew Diagnostics\n"
		out, err = runCommand("brew", "doctor")
		output += out
		if err != nil {
			return output, err
		}
		out, err = runCommand("brew", "missing")
		output += out
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

// UpdateVSCodeExt updates VSCode extensions.
func UpdateVSCodeExt() (string, error) {
	var output string
	output += "Updating VSCode Extensions\n"
	if checkCommand("code") {
		out, err := runCommand("code", "--update-extensions")
		output += out
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

// UpdateGem updates Ruby gems and clean up.
func UpdateGem() (string, error) {
	var output string
	output += "Updating Gems\n"
	gemPath, err := exec.LookPath("gem")
	if err != nil || gemPath == k_gemCmdPath {
		output += "gem is not installed."
		return output, nil
	}
	out, err := runCommand("gem", "update", "--user-install")
	output += out
	if err != nil {
		return output, err
	}
	out, err = runCommand("gem", "cleanup", "--user-install")
	output += out
	if err != nil {
		return output, err
	}
	return output, nil
}

// UpdateNodePkg updates global Node.js, npm, and Yarn packages.
func UpdateNodePkg() (string, error) {
	var output string
	output += "Updating Node Packages\n"
	if checkCommand("node") {
		output += "Updating Npm Packages\n"
		if checkCommand("npm") {
			out, err := runCommand("npm", "update", "-g")
			output += out
			if err != nil {
				return output, err
			}
		}

		output += "Updating Yarn Packages\n"
		if checkCommand("yarn") {
			out, err := runCommand("yarn", "global", "upgrade", "--latest")
			output += out
			if err != nil {
				return output, err
			}
		}
	}
	return output, nil
}

// UpdateCargo updates Rust Cargo crates by reinstalling each listed crate.
func UpdateCargo() (string, error) {
	var output string
	output += "Updating Rust Cargo Crates\n"
	if checkCommand("cargo") {
		out, _ := exec.Command("cargo", "install", "--list").Output()
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if fields := strings.Fields(line); len(fields) > 0 {
				name := fields[0]
				out, err := runCommand("cargo", "install", name)
				output += out
				if err != nil {
					return output, err
				}
			}
		}
	}
	return output, nil
}

// UpdateAppStore updates Mac App Store applications.
func UpdateAppStore() (string, error) {
	var output string
	output += "Updating App Store Applications\n"
	if checkCommand("mas") {
		out, err := runCommand("mas", "upgrade")
		output += out
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

// UpdateMacOS updates macOS system software.
func UpdateMacOS() (string, error) {
	var output string
	output += "Updating MacOS\n"
	out, err := runCommand("softwareupdate", "-i", "-a")
	output += out
	if err != nil {
		return output, err
	}
	return output, nil
}
