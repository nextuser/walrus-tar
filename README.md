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

## 使用方式

```bash
./walrus-tar -action list -blob-id ox_serTOXKIYpcn7dKXXNbJ6Y7SZWvnSN3uK71T4nig -out target
```
### 单个文件上传 下载
```bash
#上传 epochs 缺省值1 ，from指定一个文件
 ./walrus-tar -action store -from=main.go -epochs 3

# -out 写入文件，没有指定直接输出到控制台
  ./walrus-tar -action read -blob-id AjSbwCzGcZRDXRatr0aOedvKYej33w_K5wQr74KrahM -out t

```

### 单个目录上传 下载
```
ljl@ljl-i5-14400:main$ tree data/apt
data/apt
├── a
├── b
└── c
# 上传目录
 ./walrus-tar -action store -from data/apt

# 获得blobid
UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y

# 查看blob信息
 ./walrus-tar -action list -blob-id UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
url= https://aggregator.walrus-testnet.walrus.space/v1/UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
data/apt : size(0), mtime(Thu Dec 19 22:05:10 CST 2024) type(d)
data/apt/a : size(3), mtime(Thu Dec 19 23:36:01 CST 2024) type(f)
data/apt/b : size(3), mtime(Thu Dec 19 22:05:10 CST 2024) type(f)
data/apt/c : size(3), mtime(Thu Dec 19 22:05:10 CST 2024) type(f)

# 读取blob 中一个path对应的目录
ljl@ljl-i5-14400:main$ ./walrus-tar -action read -blob-id UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y  -path data/apt -out  out
url= https://aggregator.walrus-testnet.walrus.space/v1/UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
2024/12/20 00:32:53 解包： 到 out/data/apt/a ，共处理了 3 个字符的数据。
2024/12/20 00:32:53 解包： 到 out/data/apt/b ，共处理了 3 个字符的数据。
2024/12/20 00:32:53 解包： 到 out/data/apt/c ，共处理了 3 个字符的数据。

 读取blob中tar包,不指定path参数
 ./walrus-tar -action read -blob-id UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y   -out  out.tar

```