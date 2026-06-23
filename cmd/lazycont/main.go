package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	appconfig "github.com/pz/lazycont/internal/config"
	"github.com/pz/lazycont/internal/containercli"
	"github.com/pz/lazycont/internal/tui"
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		if len(args) > 1 {
			fmt.Fprintf(stderr, "lazycont: unexpected argument %q\n", args[1])
			printUsage(stderr)
			return 2
		}

		switch args[0] {
		case "--help", "-h", "help":
			printUsage(stdout)
			return 0
		case "--version", "-v", "version":
			fmt.Fprintf(stdout, "lazycont %s\n", version)
			return 0
		default:
			fmt.Fprintf(stderr, "lazycont: unexpected argument %q\n", args[0])
			printUsage(stderr)
			return 2
		}
	}

	return runTUI(stderr)
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `lazycont - terminal UI for Apple's container CLI

Usage:
  lazycont [--help] [--version]

Options:
  --help     Show this help.
  --version  Print the lazycont version.
`)
}

func runTUI(stderr io.Writer) int {
	client := containercli.New("container")
	opts := tui.Options{}
	cfg, path, err := appconfig.LoadDefault()
	if err != nil {
		opts.StartupWarning = configWarning(path, err)
	} else {
		opts.CustomCommands = customCommands(cfg.Commands)
	}
	opts.ConfigPath = path
	opts.OpenConfigCommand = openConfigCommand
	opts.LoadConfigCommands = loadConfigCommands
	program := tea.NewProgram(tui.NewWithOptions(client, opts), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(stderr, "lazycont: %v\n", err)
		return 1
	}
	return 0
}

func openConfigCommand(path string) (*exec.Cmd, error) {
	if err := appconfig.Ensure(path); err != nil {
		return nil, err
	}
	editor := strings.TrimSpace(os.Getenv("VISUAL"))
	if editor == "" {
		editor = strings.TrimSpace(os.Getenv("EDITOR"))
	}
	if editor == "" {
		editor = "vi"
	}
	return editorCommand(editor, path)
}

func editorCommand(editor string, path string) (*exec.Cmd, error) {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return nil, errors.New("editor is required")
	}
	args := append(append([]string(nil), parts[1:]...), path)
	return exec.Command(parts[0], args...), nil
}

func loadConfigCommands() ([]tui.CustomCommand, error) {
	cfg, _, err := appconfig.LoadDefault()
	if err != nil {
		return nil, err
	}
	return customCommands(cfg.Commands), nil
}

func configWarning(path string, err error) string {
	if path == "" {
		return fmt.Sprintf("config: %v", err)
	}
	return fmt.Sprintf("config %s: %v", path, err)
}

func customCommands(commands []appconfig.Command) []tui.CustomCommand {
	out := make([]tui.CustomCommand, 0, len(commands))
	for _, command := range commands {
		out = append(out, tui.CustomCommand{
			Name: command.Name,
			Args: append([]string(nil), command.Args...),
		})
	}
	return out
}
