package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"

	"github.com/fimreal/goutils/crypto"
	"github.com/fimreal/goutils/ezap"
	httpc "github.com/fimreal/goutils/http"
)

// 增
func (gc *GitClient) FileCreate(remote string, fileContent string, commitMessage string) error {
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s",
		gc.Address, url.PathEscape(gc.Project), url.PathEscape(remote))
	ezap.Debug("POST url: ", url)

	gitlabCommit := &GitlabCommit{
		Branch:        gc.Branch,
		AuthorEmail:   gc.Mail,
		AuthorName:    gc.User,
		Encoding:      "base64",
		Content:       fileContent,
		CommitMessage: commitMessage,
	}
	body, err := json.Marshal(gitlabCommit)
	if err != nil {
		return err
	}
	ezap.Debug("POST body: ", string(body))

	resp, err := httpc.HttpDo(url, "POST", body, gc.Headers)
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

// 改
func (gc *GitClient) FileUpdate(remote string, fileContent string, commitMessage string) error {
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s",
		gc.Address, url.PathEscape(gc.Project), url.PathEscape(remote))
	ezap.Debug("PUT url: ", url)

	gitlabCommit := &GitlabCommit{
		Branch:        gc.Branch,
		AuthorEmail:   gc.Mail,
		AuthorName:    gc.User,
		Encoding:      "base64",
		Content:       fileContent,
		CommitMessage: commitMessage,
	}
	body, err := json.Marshal(gitlabCommit)
	if err != nil {
		return err
	}
	ezap.Debug("PUT body: ", string(body))

	resp, err := httpc.HttpDo(url, "PUT", body, gc.Headers)
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

func (gc *GitClient) Upload(local string, remote string, commitMessage string) error {
	fileByte, err := os.ReadFile(local)
	if err != nil {
		return fmt.Errorf("读取本地文件遇到错误: %v", err)
	}
	fileContentBS64 := crypto.GetBase64Encode(string(fileByte))

	ezap.Debugf("尝试创建远程文件[%s]", remote)
	err = gc.FileCreate(remote, fileContentBS64, commitMessage)
	if err != nil {
		if fileExist, _ := regexp.Match("A file with this name already exists", []byte(err.Error())); fileExist {
			ezap.Debug(err)
			ezap.Warn("文件已存在, 更新文件内容")

			err = gc.FileUpdate(remote, fileContentBS64, commitMessage)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
