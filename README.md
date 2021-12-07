# JMS 录像上传工具

此工具主要用于 JumpServer 组件录像上传到指定存储，主要包括 OSS、Azure、OBS 和 S3 存储。

### 上传单个录像文件

例：replay_uploader [flag] [flag]

```bash
    以下是命令行参数

    --date string
    Target date, 默认当前日期，格式 2021-01-20

    --file string
    Replay file path，录像文件路径

    --key string
    Access key，组件使用的认证Key
    
    --keyfile string
    Key file，存储 access key 的文件路径
    
    --remove
    成功上传后，是否删除原 file 文件, (如果不写，默认不删除）
    
    --sid string
    Session ID, 会话ID （如果不写，则从 Replay file 文件名中提取会话 session id）

    -url string
    JumpServer URL，JMS Core的地址 (default "http://127.0.0.1:8080")

```

## 扫描目录上传遗留录像

例：replay_uploader scan  [flag] [flag]

```bash
    以下是命令行参数

    --baseDir string   需要扫描的遗留的录像文件目录

    --key string       Access key，组件使用的认证Key
    --keyfile string   Key file，存储 access key 的文件路径
    --remove           成功上传后，是否删除原 file 文件
    --url string       JumpServer URL，JMS Core的地址 (default "http://127.0.0.1:8080")

```

### 参数注意事项

key 的值是 base64 之后的结果，例如通过查看组件的key：`9b918f99-29b5-4e99-8c4d-aaba19271111:7568b6da-9d37-46cb-8f4f-cab2e196db23`，
然后 base64 得到 `OWI5MThmOTktMjliNS00ZTk5LThjNGQtYWFiYTE5MjcxMTExOjc1NjhiNmRhLTlkMzctNDZjYi04ZjRmLWNhYjJlMTk2ZGIyMwo=`

keyfile 是存储`access key`的文件路径，内容格式类似 `9b918f99-29b5-4e99-8c4d-aaba19271111:7568b6da-9d37-46cb-8f4f-cab2e196db23`