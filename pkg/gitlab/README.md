## Gitlab 相关 API 参考

https://docs.gitlab.com/ee/api/

### 本地运行 Gitlab 帮助文档：

```bash
docker run -it --rm -p 4000:4000 registry.gitlab.com/gitlab-org/gitlab-docs:<gitlab version>
# docker run -it --rm -p 4000:4000 registry.gitlab.com/gitlab-org/gitlab-docs:11.11
```

#### 查看文件内容

```bash
curl -H "PRIVATE-TOKEN: <gitlab token>" http://<gitlab host>/api/v4/projects/<repo id>/repository/files/<floder>%2F<filename>/raw?ref=<branch|tag|commitid>
```

#### 创建新文件

```bash
curl -XPOST -H "PRIVATE-TOKEN: <gitlab token>" \
-H "Content-Type: application/json" \-d '{
    "branch": "master",
    "author_email": "xm@epurs.com",
    "author_name": "xm",
    "encoding": "base64",
    "content": "YnM2NCBjb250ZW50Cg==",
    "commit_message": "use api to create a new file with curl"
}' \
http://<gitlab host>/api/v4/projects/<repo id>/repository/files/<floder>%2F<filename.xxx>

# result: {"file_path":"xxx/xxx.xxx","branch":"xxx"} or {"message":"A file with this name already exists"}
```

#### 更新文件

同上创建文件，把 Method 改为 PUT，返回同上

```bash
curl -XPUT -H "PRIVATE-TOKEN: <gitlab token>" \
-H "Content-Type: application/json" \-d '{
    "branch": "master",
    "author_email": "xm@epurs.com",
    "author_name": "xm",
    "encoding": "base64",
    "content": "YnM2NCBjb250ZW50Cg==",
    "commit_message": "use api to update a new file with curl"
}' \
http://<gitlab host>/api/v4/projects/<repo id>/repository/files/<floder>%2F<filename.xxx>
```

#### 删除文件

```bash
curl -XDELETE -H "PRIVATE-TOKEN: <gitlab token>" \
-H "Content-Type: application/json" \-d '{
    "branch": "master",
    "author_email": "xm@epurs.com",
    "author_name": "xm",
    "commit_message": "use api to delete a new file with curl"
}' \
http://<gitlab host>/api/v4/projects/<repo id>/repository/files/<floder>%2F<filename.xxx>
```