package gitlab

import (
	"os"
	"testing"

	"github.com/fimreal/goutils/ezap"
)

var (
	gcc           string
	commitMessage = "commit by gitops-tools, please use --commit-message set this message"
)

func init() {
	var b []byte
	b, err := os.ReadFile(".token")
	if err != nil {
		os.Exit(1)
	}
	gcc = string(b)
	ezap.SetLevel("debug")
	ezap.SetLogTime("")
}

func TestParseGcc(t *testing.T) {
	ezap.SetLevel("debug")

	SetProvider("gitlab")
	gcc = "http://lxm:token@www.gitlab.cn/lxm/project-1?branch=master&mail=lxm@mail.com"
	_, err := ParseGcc(gcc)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(RemoteFilename)
}
