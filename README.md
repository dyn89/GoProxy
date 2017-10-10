# GoProxy 反向代理:中间人接力

程序自包含, 只使用到官方库

# 用户手册

无认证模式:

```
proxy_amd64 --remote=http://www.zhihu.com --local=127.0.0.1:8889--auth=no
```

亚马逊认证模式:

```
proxy_amd64 --access-key=1213 --secret-key=131312 --remote=https://search-asdsaddsa.us-east-1.es.amazonaws.com  --local=127.0.0.1:8888
```

参数:

```
  -access-key string
        公钥
  -auth string
        认证模式: aws-es|no (default "aws-es")
  -awsregion string
        aws区域 (default "us-east-1")
  -local string
        本地监听 (default "0.0.0.0:8888")
  -remote string
        代理网站
  -secret-key string
        秘钥
```

# 编译

安装Golang环境

```
./build.sh
```
