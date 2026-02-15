package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codemint/codemint-cli/internal/manifest"
	"github.com/codemint/codemint-cli/internal/tooling"
)

func resolveAITool(store *manifest.Store, override string, nonInteractive bool) (string, error) {
	if override != "" {
		if err := tooling.Validate(override); err != nil {
			return "", err
		}
		settings, err := store.LoadSettings()
		if err != nil {
			return "", err
		}
		settings.AITool = override
		if err := store.SaveSettings(settings); err != nil {
			return "", err
		}
		return override, nil
	}

	settings, err := store.LoadSettings()
	if err != nil {
		return "", err
	}
	if settings.AITool != "" {
		if err := tooling.Validate(settings.AITool); err == nil {
			return settings.AITool, nil
		}
	}
	if nonInteractive {
		return "", fmt.Errorf("no AI tool selected. run add with --tool <name> (supported: %s)", strings.Join(tooling.Supported(), ", "))
	}
	tool, err := promptAITool()
	if err != nil {
		return "", err
	}
	settings.AITool = tool
	if err := store.SaveSettings(settings); err != nil {
		return "", err
	}
	return tool, nil
}

func promptAITool() (string, error) {
	choices := tooling.Supported()
	fmt.Println("Select the AI coding tool you use for this repo:")
	for i, c := range choices {
		fmt.Printf("  %d) %s\n", i+1, c)
	}
	fmt.Print("Enter number or tool name: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" {
		return "", fmt.Errorf("no selection provided")
	}
	if idx, err := strconv.Atoi(line); err == nil {
		if idx < 1 || idx > len(choices) {
			return "", fmt.Errorf("invalid selection %d", idx)
		}
		return choices[idx-1], nil
	}
	if err := tooling.Validate(line); err != nil {
		return "", err
	}
	return line, nil
}
