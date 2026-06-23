package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadMissingConfigReturnsEmptyConfig(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Commands) != 0 {
		t.Fatalf("commands = %#v, want empty", cfg.Commands)
	}
}

func TestLoadCustomCommands(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	writeConfig(t, path, `{
		"commands": [
			{"name": " Images ", "args": [" image ", " list ", "--format", " json "]},
			{"name": "Disk usage", "args": ["system", "df"]},
			{"name": "Empty env", "args": ["run", "--env", ""]}
		]
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	want := []Command{
		{Name: "Images", Args: []string{"image", "list", "--format", "json"}},
		{Name: "Disk usage", Args: []string{"system", "df"}},
		{Name: "Empty env", Args: []string{"run", "--env", ""}},
	}
	if !reflect.DeepEqual(cfg.Commands, want) {
		t.Fatalf("commands mismatch\nwant: %#v\n got: %#v", want, cfg.Commands)
	}
}

func TestLoadRejectsInvalidCustomCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	writeConfig(t, path, `{"commands": [{"name": "Broken", "args": []}]}`)

	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "commands[0].args") {
		t.Fatalf("err = %v, want args validation error", err)
	}
}

func TestLoadRejectsTrailingJSONData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	writeConfig(t, path, `{"commands": []} {"commands": []}`)

	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "trailing JSON data") {
		t.Fatalf("err = %v, want trailing data error", err)
	}
}

func writeConfig(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
