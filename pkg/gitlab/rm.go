package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/fimreal/goutils/ezap"
	httpc "github.com/fimreal/goutils/http"
)

// 删
func (gc *GitClient) FileDelete(remote string, commitMessage string) error {

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s",
		gc.Address, url.PathEscape(gc.Project), url.PathEscape(remote))
	ezap.Debug("DELETE url: ", url)
	gitlabCommit := &GitlabCommit{
		Branch:        gc.Branch,
		AuthorEmail:   gc.Mail,
		AuthorName:    gc.User,
		CommitMessage: commitMessage,
	}

	body, err := json.Marshal(gitlabCommit)
	if err != nil {
		return err
	}
	ezap.Debug("DELETE body: ", string(body))

	resp, err := httpc.HttpDo(url, "DELETE", body, gc.Headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	ezap.Debug(string(respBody))
	if err != nil {
		return err
	} else if resp.StatusCode >= 400 {
		return fmt.Errorf("%d, %s", resp.StatusCode, string(respBody))
	}

	return nil
}
