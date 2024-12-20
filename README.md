# 简介
- 用于将目录打包成tar包，上传成walrus的一个blob
- 从walrus 获取一个blob，提取tar包一个文件或目录
# 开发说明
## 编译
``
go env -w GO111MODULE=off
cd main 
go build -o ../release/wtar

``

## 单元测试代码

```test
go env -w GO111MODULE=off
cd main
go test

```

# 使用说明
- 命令
```bash
./wtar -action list -blob-id ox_serTOXKIYpcn7dKXXNbJ6Y7SZWvnSN3uK71T4nig -out target
```
### 单个文件上传 下载

#上传 epochs 缺省值1 ，from指定一个文件
- 命令
```bash
./wtar -c store -f=main.go -e 3
```
- 输出结果
```
url= https://publisher.walrus-testnet.walrus.space/v1/store?epochs=3
2024/12/20 09:06:44 成功打包 main.go ，共写入了 5006 字节的数据
get response:
{"newlyCreated":{"blobObject":{"id":"0x11197c633d117e8ea56f5aaa9b3e2f5323a835532b80cd4a8089ffc71f6a7c53","registeredEpoch":63,"blobId":"QXMOLMtbj428L_Bgc7To2me5g7MYmdvw8mGDAMlk0Gs","size":6656,"encodingType":"RedStuff","certifiedEpoch":63,"storage":{"id":"0x75f7617d2adec8b156a6a64f9ac2f8bcbee1f9dd070593d2051cf53d6fa5cf1d","startEpoch":63,"endEpoch":66,"storageSize":65023000},"delet
```

# 列出walrus上tar文件的路径
- 命令
```bash
 ./wtar -c list -b QXMOLMtbj428L_Bgc7To2me5g7MYmdvw8mGDAMlk0Gs
```
- 输出结果
```
url= https://aggregator.walrus-testnet.walrus.space/v1/QXMOLMtbj428L_Bgc7To2me5g7MYmdvw8mGDAMlk0Gs
------------------list file path in tar---------------
main.go : size(5006), mtime(Fri Dec 20 08:50:54 CST 2024) type(f)
```

# 读取上传单个文件内容
- 命令
```bash
./wtar -c read -b QXMOLMtbj428L_Bgc7To2me5g7MYmdvw8mGDAMlk0Gs -p main.go
```

### 单个目录上传 下载
####  查看一个目录结构
- 命令
```bash
$ tree data/apt
```
- 输出结果
```
data/apt
├── a
├── b
└── c
```
#### 上传目录
- 命令
```bash
./wtar  -c store -f data/apt
```
- 输出结果
```
url= https://publisher.walrus-testnet.walrus.space/v1/store?epochs=1
2024/12/20 08:59:15 成功打包 data/apt/a ，共写入了 3 字节的数据
2024/12/20 08:59:15 成功打包 data/apt/b ，共写入了 3 字节的数据
2024/12/20 08:59:15 成功打包 data/apt/c ，共写入了 3 字节的数据

get response:
{"alreadyCertified":{"blobId":"UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y","eventOrObject":{"Event":{"txDigest":"77HyHGjkJ1gKhriqWBFLDPJTNXNeExntwA4vRgi4VaWf","eventSeq":"0"}},"endEpoch":64}}
```
- 获得blobid
UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y

#### 列出tar包文件中的文件路径列表
- 命令
```bash
./wtar -c list -b UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
```
- 输出结果
```
url= https://aggregator.walrus-testnet.walrus.space/v1/UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
------------------list file path in tar---------------
data/apt : size(0), mtime(Thu Dec 19 22:05:10 CST 2024) type(d)
data/apt/a : size(3), mtime(Thu Dec 19 23:36:01 CST 2024) type(f)
data/apt/b : size(3), mtime(Thu Dec 19 22:05:10 CST 2024) type(f)
data/apt/c : size(3), mtime(Thu Dec 19 22:05:10 CST 2024) type(f)
```


#### 读取blob 中一个文件路径
- 命令
```bash
$ ./wtar  -c read -b UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y -p data/apt/a
```
- 输出结果
```
url= https://aggregator.walrus-testnet.walrus.space/v1/UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
aa
--------解包： 到  ，共处理了 3 个字符的数据。--------
$ ./wtar  -c read -b UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y -p data/apt/c
url= https://aggregator.walrus-testnet.walrus.space/v1/UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
cc
--------解包： 到  ，共处理了 3 个字符的数据。--------
```

#### 读取blob 中一个path对应的目录，写入result目录下
- 命令
```bash
$ ./wtar  -c read -b UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y -p data/apt -o  result
```
- 输出结果
```
url= https://aggregator.walrus-testnet.walrus.space/v1/UhdFYLm5Qf3sKdSNno5XgieVOUQYk0fqEVB3rqiKl_Y
--------解包： 到 result/a ，共处理了 3 个字符的数据。--------
--------解包： 到 result/b ，共处理了 3 个字符的数据。--------
--------解包： 到 result/c ，共处理了 3 个字符的数据。--------
$ tree result
result
├── a
├── b
└── c


```