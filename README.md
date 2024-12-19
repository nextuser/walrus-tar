# walrus-tar
- 用于将目录打包成tar包，上传成walrus的一个blob
- 从walrus 获取一个blob，提取tar包一个文件或目录
## main
``
go env -w GO111MODULE=off
cd main 
go build -o ../release/walrus-tar
./walrus-tar
``

## tar file

```test
go env -w GO111MODULE=off
cd main
go test

```
