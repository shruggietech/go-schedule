package clievaluate

import (
	"bytes"
	"embed"
	"errors"
	"fmt"

	"github.com/shruggietech/go-schedule/internal/cli"
)

//go:embed inventory-v1.json
var inventoryFiles embed.FS

// ActualCommands reads the actual initialized Cobra metadata, without executing it.
func ActualCommands() ([]Command, error) {
	data, err := cli.CommandInventoryJSON()
	if err != nil {
		return nil, err
	}
	var commands []Command
	if err := Decode(bytes.NewReader(data), &commands); err != nil {
		return nil, err
	}
	return commands, nil
}

// CheckInventory requires deliberate review of any changed command/flag metadata.
func CheckInventory(commands []Command) error {
	want, err := inventoryFiles.ReadFile("inventory-v1.json")
	if err != nil {
		return err
	}
	got, err := Encode(commands)
	if err != nil {
		return err
	}
	if !bytes.Equal(want, got) {
		return errors.New("CLI inventory drift: review actual command metadata and update versioned scenarios/inventory before claiming coverage")
	}
	return nil
}

// InventoryCoverage checks required class bindings without counting help as behavior.
func InventoryCoverage(m Manifest) error {
	if err := CheckInventory(m.Commands); err != nil {
		return err
	}
	if err := ValidateManifest(m); err != nil {
		return err
	}
	covered := map[string]bool{}
	for _, c := range m.Cases {
		for _, item := range c.Coverage {
			covered[item] = true
		}
	}
	for _, c := range m.Commands {
		for _, flag := range c.Flags {
			key := c.Path + " --" + flag.Name
			if flag.Scope == "inherited" {
				key = " --" + flag.Name
			}
			if !covered[key] {
				return fmt.Errorf("uncovered effective flag: %s", key)
			}
		}
		for _, alias := range c.Aliases {
			if !covered[c.Path+" alias="+alias] {
				return fmt.Errorf("uncovered command alias: %s", alias)
			}
		}
	}
	return nil
}
