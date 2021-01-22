# JMS Upload tool

此工具主要用于 JumpServer 组件录像上传到指定存储，主要包括 OSS、Azure 和 S3 存储。


```
    以下是命令行参数

    -date string
    Target date, 默认当前日期，格式 2021-01-20

    -file string
    Replay file path，录像文件路径

    -key string
    Access key，组件使用的认证Key

    -url string
    JumpServer URL，JMS Core的地址 (default "http://127.0.0.1:8080")

```
## 参数说明

key 的值是 base64之后的结果，例如通过查看组件的key：`9b918f99-29b5-4e99-8c4d-aaba19271111:7568b6da-9d37-46cb-8f4f-cab2e196db23`，
然后 base64 得到 `OWI5MThmOTktMjliNS00ZTk5LThjNGQtYWFiYTE5MjcxMTExOjc1NjhiNmRhLTlkMzctNDZjYi04ZjRmLWNhYjJlMTk2ZGIyMwo=`

