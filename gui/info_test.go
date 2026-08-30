package gui

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func TestEmbeddedBrandIconsUseApprovedSizes(t *testing.T) {
	for _, test := range []struct {
		name   string
		data   []byte
		width  int
		height int
	}{
		{name: "full mark", data: iconBytes, width: 1024, height: 1024},
		{name: "reduced mark", data: windowIconBytes, width: 256, height: 256},
	} {
		t.Run(test.name, func(t *testing.T) {
			config, _, err := image.DecodeConfig(bytes.NewReader(test.data))
			if err != nil {
				t.Fatalf("decode embedded icon: %v", err)
			}
			if config.Width != test.width || config.Height != test.height {
				t.Fatalf("embedded icon is %dx%d, want %dx%d", config.Width, config.Height, test.width, test.height)
			}
		})
	}
}

func TestUI_InfoContentUsesLocalIdentityAndExactLinks(t *testing.T) {
	const version = "2026.08.28-development-build-with-long-id"
	content := buildInfoContent(version)

	var images []*canvas.Image
	var labels []string
	links := make(map[string]string)
	walkInfoObjects(content, func(object fyne.CanvasObject) {
		switch object := object.(type) {
		case *canvas.Image:
			images = append(images, object)
		case *widget.Label:
			labels = append(labels, object.Text)
		case *widget.Hyperlink:
			links[object.Text] = object.URL.String()
		}
	})

	if len(images) != 1 || images[0].Resource != appIcon {
		t.Fatalf("Info image does not reuse appIcon: %+v", images)
	}
	if images[0].FillMode != canvas.ImageFillContain {
		t.Fatalf("Info image fill mode = %v, want contain", images[0].FillMode)
	}
	assertInfoLabel(t, labels, "Version "+version)
	assertInfoLabel(t, labels, "Built and maintained by ShruggieTech")

	wantLinks := map[string]string{
		"ShruggieTech":      "https://shruggie.tech",
		"Source repository": "https://github.com/shruggietech/go-schedule",
		"Documentation":     "https://shruggietech.github.io/go-schedule/",
	}
	for text, wantURL := range wantLinks {
		if got := links[text]; got != wantURL {
			t.Errorf("Info link %q = %q, want %q", text, got, wantURL)
		}
	}
}

func walkInfoObjects(object fyne.CanvasObject, visit func(fyne.CanvasObject)) {
	visit(object)
	switch object := object.(type) {
	case *fyne.Container:
		for _, child := range object.Objects {
			walkInfoObjects(child, visit)
		}
	case *container.Scroll:
		walkInfoObjects(object.Content, visit)
	}
}

func assertInfoLabel(t *testing.T, labels []string, want string) {
	t.Helper()
	for _, label := range labels {
		if label == want {
			return
		}
	}
	t.Errorf("Info labels %q do not contain %q", labels, want)
}
