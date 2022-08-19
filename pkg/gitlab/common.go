package gitlab

var (
	gitClient *GitClient
)

type GitClient struct {
	Provider string
	// Scheme   string // 并入 address
	User string
	// Token   string //并入 headers
	Address string
	Project string
	Branch  string
	Mail    string
	Headers map[string]string
}

type GitlabCommit struct {
	// 分支，例如 master、devlop
	Branch string `json:"branch"`
	// 修改人邮箱，例如 lxm@epurs.com
	AuthorEmail string `json:"author_email"`
	// 修改人名字，例如 lxm
	AuthorName string `json:"author_name"`
	// 内容编码，使用 base64，删除时无用
	Encoding string `json:"encoding"`
	// 文件内容，删除时无用
	Content string `json:"content"`
	// 提交时备注信息
	CommitMessage string `json:"commit_message"`
}

// type GitlabFile struct {
// 	Filename    string
// 	ProjectName string
// 	ref         string
// 	RawContent  []byte
// }
