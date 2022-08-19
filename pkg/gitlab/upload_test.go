package gitlab

import "testing"

func TestFileCreate(t *testing.T) {
	gc, err := ParseGcc(gcc)
	if err != nil {
		t.Fatal(err)
	}
	err = gc.FileCreate("testcreate1", "abcdefg", commitMessage)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileUpdate(t *testing.T) {
	gc, err := ParseGcc(gcc)
	if err != nil {
		t.Fatal(err)
	}
	err = gc.FileUpdate("testcreate1", "abcdefghijklmn", commitMessage)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpload(t *testing.T) {
	gc, err := ParseGcc(gcc)
	if err != nil {
		t.Fatal(err)
	}
	err = gc.Upload("README.md", "testupload.md", commitMessage)
	if err != nil {
		t.Fatal(err)
	}
}
