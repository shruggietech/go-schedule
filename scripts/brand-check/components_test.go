package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCanonicalFormWrappersPreserveBaseClasses(t *testing.T) {
	for _, test := range []struct {
		path      string
		baseClass string
	}{
		{path: "Input.jsx", baseClass: "gs-input"},
		{path: "Select.jsx", baseClass: "gs-select"},
	} {
		t.Run(test.path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "..", "brand", "components", "forms", test.path))
			if err != nil {
				t.Fatal(err)
			}
			source := string(data)
			if !strings.Contains(source, `className = ""`) ||
				!strings.Contains(source, "`"+test.baseClass+` ${className}`+"`") {
				t.Fatalf("%s must merge a caller className after %s", test.path, test.baseClass)
			}
		})
	}
}

func TestCanonicalLightButtonsUseAccessibleForeground(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "brand", "components", "components.css"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), ".gs-light .gs-button { color: #FFFFFF; }") {
		t.Fatal("light-surface primary buttons must use a white foreground")
	}
}
