# Reverse Proxy: Use For AWS IAM/v4 Auth

The program itself contains only the official library

# How to Use?

1.No authentication proxy way:

```
proxy_amd64 --remote=http://www.zhihu.com --local=127.0.0.1:8889--auth=no
```

2.AWS authentication proxy way:

```
proxy_amd64 --access-key=1213 --secret-key=131312 --remote=https://search-asdsaddsa.us-east-1.es.amazonaws.com  --local=127.0.0.1:8888
```

Then you can have a try:

```
curl 127.0.0.1:8888
```

If you want to write some other aws service auth you can refer core code:

```
	.....
	
	switch *Type {
	case "aws-es":
		AwsConfig.AwsService = "es"    //   Just Modify here
		AwsConfig.AwsRegion = *AWSRegion  // Just Modify here
	default:
		break
	}
	AwsConfig.AwsID = *Key
	AwsConfig.AwsKey = *Secret
   
	....
	
	req.ParseForm()
    amzdate, authorization_header := AwsAuthSignature(AwsConfig, getURIPath(req.URL), req.Method, req.URL.Host, req.Form, buf)
    req.Header.Set("X-Amz-Date", amzdate)
    req.Header.Set("Authorization", authorization_header)

```

Parm:

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

# How to Make Executable File?

Install Golang and `git clone https://github.com/hunterhug/GoProxy`

Then do this:

```
./build.sh
```
