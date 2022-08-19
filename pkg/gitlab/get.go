package gitlab

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/fimreal/goutils/ezap"
	httpc "github.com/fimreal/goutils/http"
)

// 获取文件详细信息，内容 base64
//
// return string
func (gc *GitClient) GetFileInfo(file string) (string, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s?ref=%s",
		gc.Address, url.PathEscape(gc.Project), url.PathEscape(file), gc.Branch)
	ezap.Debug("Get url: ", url)

	resp, err := httpc.HttpDo(url, "GET", nil, gc.Headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fileInfo, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("请求出错[%d], Err: %v, Response Message: %s", resp.StatusCode, err, string(fileInfo))
	}

	return string(fileInfo), nil
}

// 查
// return []byte
func (gc *GitClient) GetFileRaw(remote string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s/raw?ref=%s",
		gc.Address, url.PathEscape(gc.Project), url.PathEscape(remote), gc.Branch)
	ezap.Debug("GET url: ", url)

	resp, err := httpc.HttpDo(url, "GET", nil, gc.Headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileByte, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("请求出错[%d], Err: %v, Response Message: %s", resp.StatusCode, err, string(fileByte))
	}

	return fileByte, nil
}

func (gc *GitClient) Download(remote string, local string) (err error) {
	var fileContent []byte
	if fileContent, err = gc.GetFileRaw(remote); err != nil {
		return fmt.Errorf("获取远程文件出错， %v", err)
	}

	file, err := os.Create(local)
	if err != nil {
		return fmt.Errorf("创建本地文件遇到错误: %v", err)
	}
	defer file.Close()

	_, err = file.Write(fileContent)
	if err != nil {
		return fmt.Errorf("写入本地文件出错， %v", err)
	}
	return nil
}
