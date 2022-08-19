package gitlab

import (
	"testing"
)

func TestGetFileInfo(t *testing.T) {
	gc, err := ParseGcc(gcc)
	if err != nil {
		t.Fatal(err)
	}
	info, err := gc.GetFileInfo("README.md")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}

func TestGetFileRaw(t *testing.T) {
	gc, err := ParseGcc(gcc)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := gc.GetFileRaw("README.md")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(raw))
}

func TestDownload(t *testing.T) {
	gc, err := ParseGcc(gcc)
	if err != nil {
		t.Fatal(err)
	}
	err = gc.Download("README.md", "download.md")
	if err != nil {
		t.Fatal(err)
	}
}
