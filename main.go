package main

import (
	"bytes"
	"flag"
	"fmt"
	"macup/macup"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Entry point of the update script
func main() {
	cli := flag.Bool("cli", false, "Run in command-line mode")
	flag.Parse()

	if *cli {
		runCLI()
	} else {
		runGUI()
	}
}

func runGUI() {
	a := app.New()
	w := a.NewWindow("macup GUI")
	w.Resize(fyne.NewSize(500, 800)) // Set window size to 600x500

	output := widget.NewMultiLineEntry()
	output.SetMinRowsVisible(10)

	var selectedUpdates []string
	updateChecks := []*widget.Check{}
	for _, u := range macup.Updates {
		update := u
		check := widget.NewCheck(update.Name, func(checked bool) {
			if checked {
				selectedUpdates = append(selectedUpdates, update.Name)
			} else {
				// remove from selected updates
				for i, v := range selectedUpdates {
					if v == update.Name {
						selectedUpdates = append(selectedUpdates[:i], selectedUpdates[i+1:]...)
						break
					}
				}
			}
		})
		updateChecks = append(updateChecks, check)
	}

	updateButton := widget.NewButton("Update", func() {
		output.SetText("Running updates...")
		var outputText string
		for _, updateName := range selectedUpdates {
			for _, u := range macup.Updates {
				if u.Name == updateName {
					//outputText += fmt.Sprintf("---\nRunning %s ---\n", u.Name)
					cmdOutput, err := u.Run()
					if err != nil {
						outputText += fmt.Sprintf("Error: %v\n", err)
					}
					outputText += cmdOutput
					outputText += "\n"
					output.SetText(outputText)
				}
			}
		}
		output.SetText(outputText + "--- Updates finished ---")
	})

	checkContainer := container.NewVBox()
	for _, check := range updateChecks {
		checkContainer.Add(check)
	}

	w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Select updates to install:"),
			checkContainer,
			updateButton,
		), nil, nil, nil,
		container.NewScroll(output),
	))

	w.ShowAndRun()
}

// Entry point of the update script
func runCLI() {
	// Define and parse the --yes flag
	yes := flag.Bool("yes", false, "Use previous selections without prompting")
	flag.Parse()

	// Check for internet connectivity before proceeding
	if !macup.CheckInternet() {
		os.Exit(1)
	}

	// Load user's previous selections
	config, err := macup.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Prompt the user to select updates
	var selectedUpdates []string
	if *yes && len(config.SelectedUpdates) > 0 {
		selectedUpdates = config.SelectedUpdates
	} else {
		if *yes {
			println("No previous updates selected. Prompting for selection.")
		}
		var err error
		selectedUpdates, err = macup.SelectUpdates(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting updates: %v\n", err)
			os.Exit(1)
		}

		// Save the user's selections
		config.SelectedUpdates = selectedUpdates
		if err := config.SaveConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
	}

	// Run the selected updates
	buffers := make(map[string]*bytes.Buffer)
	for _, updateName := range selectedUpdates {
		buffers[updateName] = &bytes.Buffer{}
	}
	for _, updateName := range selectedUpdates {
		buffers[updateName].WriteString(fmt.Sprintf("---\nRunning %s ---\n", updateName))
		for _, u := range macup.Updates {
			if u.Name == updateName {
				cmdOutput, err := u.Run()
				if err != nil {
					buffers[updateName].WriteString(fmt.Sprintf("Error: %v\n", err))
				}
				buffers[updateName].WriteString(cmdOutput)
				buffers[updateName].WriteString("\n")
			}
		}
	}

	// Print the output of each update function
	for _, updateName := range selectedUpdates {
		if buffer, ok := buffers[updateName]; ok {
			fmt.Println(buffer.String())
		}
	}
}
