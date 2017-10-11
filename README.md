# Reverse Proxy: Use For AWS IAM/v4 Auth

The program itself contains only the official library

# How to Use?

1.No authentication proxy way:

```
proxyamd64 --remote=http://www.zhihu.com --local=127.0.0.1:8889 --auth=no
```

2.AWS authentication proxy way:

```
proxyamd64 --access-key=1213 --secret-key=131312 --remote=https://search-asdsaddsa.us-east-1.es.amazonaws.com  --local=127.0.0.1:8888
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
	
	.....

	req.ParseForm()
	amzdate, authorization_header := AwsAuthSignature(AwsConfig, getURIPath(req.URL), req.Method, req.URL.Host, req.Form, buf)
	req.Header.Set("X-Amz-Date", amzdate)
	req.Header.Set("Authorization", authorization_header)
	
	.....

```

Usage of proxyamd64:

```
  -access-key string
        access key
  -auth string
        auth way: aws-es|no (default "aws-es")
  -aws-region string
        aws region(onlu valid in aws auth way) (default "us-east-1")
  -local string
        local proxy address (default "0.0.0.0:8888")
  -remote string
        remote web such as http://www.google.com(must have http)
  -secret-key string
        secret key
```

# How to Make Executable File?

Install Golang and `git clone https://github.com/hunterhug/GoProxy`

Then do this:

```
./build.sh
```

You can download in [https://github.com/hunterhug/GoProxy/releases](https://github.com/hunterhug/GoProxy/releases)