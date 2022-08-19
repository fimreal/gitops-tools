package gitlab

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fimreal/goutils/ezap"
)

func init() {
	gitClient = &GitClient{}
}

func SetProvider(p string) {
	gitClient.Provider = p
}

func ParseGcc(gcc string) (*GitClient, error) {
	u, err := url.Parse(gcc)
	if err != nil {
		return nil, err
	}
	// if only token privide, token will set to user
	user := u.User.Username()
	token, found := u.User.Password()
	if !found {
		return gitClient, fmt.Errorf("not found user/token filed")
	}

	queryMap, _ := url.ParseQuery(u.RawQuery)
	// 默认 master 分支
	branch := "master"
	// 默认 mail 地址
	mail := "gitops@mail.com"
	// 创建 headers
	headers := make(map[string]string)
	headers["PRIVATE-TOKEN"] = token
	headers["User-Agent"] = "gohttp"
	headers["Accept"] = "*/*"
	headers["Content-Type"] = "application/json"

	for k, v := range queryMap {
		switch k {
		case "branch":
			branch = v[0]
		case "mail":
			mail = v[0]
		case "useragent":
			headers["User-Agent"] = v[0]
		default:
			ezap.Debugf("发现未知参数: %s:%v", k, v)
		}
	}

	gitClient = &GitClient{
		Provider: gitClient.Provider,
		User:     user,
		Address:  u.Scheme + "://" + u.Host,
		Project:  strings.TrimPrefix(u.Path, "/"),
		Branch:   branch,
		Mail:     mail,
		Headers:  headers,
	}

	ezap.Debugf("%+v\n", *gitClient)
	return gitClient, nil
}
