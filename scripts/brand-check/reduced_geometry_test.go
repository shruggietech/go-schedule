package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func reducedGeometry(t *testing.T, filename string) []string {
	t.Helper()
	path := filepath.Join("..", "..", "brand", "logos", "svg", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var shapes []string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return shapes
		}
		if err != nil {
			t.Fatal(err)
		}
		element, ok := token.(xml.StartElement)
		if !ok || (element.Name.Local != "path" && element.Name.Local != "rect") {
			continue
		}
		var attributes []string
		for _, attr := range element.Attr {
			if attr.Name.Local != "fill" {
				attributes = append(attributes, attr.Name.Local+"="+attr.Value)
			}
		}
		shapes = append(shapes, element.Name.Local+":"+strings.Join(attributes, ";"))
	}
}

func TestMonochromeReducedMarksPreserveCanonicalGeometry(t *testing.T) {
	t.Parallel()
	canonical := reducedGeometry(t, "go-schedule-mark-reduced.svg")
	if len(canonical) == 0 {
		t.Fatal("canonical reduced mark has no geometry")
	}
	for _, filename := range []string{"go-schedule-mark-reduced-white.svg", "go-schedule-mark-reduced-black.svg"} {
		if got := reducedGeometry(t, filename); !reflect.DeepEqual(got, canonical) {
			t.Errorf("%s changed canonical reduced geometry", filename)
		}
	}
}
