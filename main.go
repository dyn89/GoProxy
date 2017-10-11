/*
   Created by jinhan on 17-10-10.
   Tip:  Reverse Proxy: Use For AWS IAM/v4 Auth
   Update:
*/
package main

import (
	//"errors"
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type handle struct {
	Https bool
	host  string
	port  string
	local string
}

var (
	Version   = "v1.1"
	Key       = flag.String("access-key", "", "access key")
	Secret    = flag.String("secret-key", "", "secret key")
	Remote    = flag.String("remote", "", "remote web such as http://www.google.com(must have http)")
	Local     = flag.String("local", "0.0.0.0:8888", "local proxy address")
	Type      = flag.String("auth", "aws-es", "auth way: aws-es|no")
	AWSRegion = flag.String("aws-region", "us-east-1", "aws region(onlu valid in aws auth way)")
	AwsConfig = AwsAuth{}
	Proxy     *handle
)

func init() {
	flag.Parse()
	if *Remote == "" {
		fmt.Println("remote empty")
		os.Exit(1)
	}

	switch *Type {
	case "aws-es":
		AwsConfig.AwsService = "es"
		AwsConfig.AwsRegion = *AWSRegion
	default:
		break
	}
	AwsConfig.AwsID = *Key
	AwsConfig.AwsKey = *Secret
}

func main() {
	fmt.Printf("Reverse Proxy %s Start\nRep: https://github.com/hunterhug/GoProxy \n------------------\n", Version)
	startServer()
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	st := "http://" + this.host
	if this.Https {
		st = "https://" + this.host
	}
	if this.port != "80" {
		st = st + ":" + this.port
	}
	remote, err := url.Parse(st)
	if err != nil {
		panic(err)
	}
	proxy := NewAWSReverseProxy(remote)
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   50 * time.Second,
			KeepAlive: 50 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   50 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
	}
	proxy.ServeHTTP(w, r)
}

func startServer() {
	phttp := strings.Split(*Remote, "//")
	temp := strings.Split(phttp[1], "/")[0]
	h := strings.Split(temp, ":")
	if len(h) == 1 {
		Proxy = &handle{host: h[0], port: "80", local: *Local, Https: false}
	} else {
		Proxy = &handle{host: h[0], port: h[1], local: *Local, Https: false}
	}
	if strings.Contains(phttp[0], "https") {
		Proxy.Https = true
	}

	fmt.Printf("Use HTTPS: %v\n", Proxy.Https)
	fmt.Printf("Remote WEB: %v\n", Proxy.host+":"+Proxy.port)
	fmt.Printf("Local Proxy: %v\n", Proxy.local)
	fmt.Printf("------------------\njust curl %v\n------------------\n", Proxy.local)
	err := http.ListenAndServe(Proxy.local, Proxy)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func refplace(ref string) string {
	return strings.Replace(ref, strings.Split(strings.Split(ref, "//")[1], "/")[0], Proxy.host, -1)
}

func locationplace(ref string) string {
	return strings.Replace(ref, strings.Split(strings.Split(ref, "//")[1], "/")[0], Proxy.local, -1)
}

func NewAWSReverseProxy(target *url.URL) *httputil.ReverseProxy {
	fmt.Println("**********")
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}

		fmt.Printf("URL: %#v\n", req.URL)

		ref := req.Header.Get("Referer")

		if ref != "" {
			req.Header.Set("Referer", refplace(ref))
		}

		cl := req.Header.Get("Content-Length")

		buf := []byte(nil)
		if cl != "" && cl != "0" {
			buf, _ = ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

			fmt.Printf("data: %#v\n", string(buf))
		}

		switch *Type {
		case "aws-es":
			// 在此验证
			req.ParseForm()
			//fmt.Println(AwsConfig, req.URL.Path, req.Method, req.URL.Host, req.Form)
			amzdate, authorization_header := AwsAuthSignature(AwsConfig, UriEncode(req.URL.Path, true), req.Method, req.URL.Host, req.Form, buf)
			req.Header.Set("X-Amz-Date", amzdate)
			req.Header.Set("Authorization", authorization_header)
		default:
			break
		}

		req.Header.Set("Host", req.URL.Host)
		req.Host = req.URL.Host

		for k, v := range req.Header {
			fmt.Println(k, v)
		}
		fmt.Println("-----------------")
	}

	modify := func(rsp *http.Response) error {

		ct := rsp.Header.Get("Content-Type")
		//return errors.New("diy error")
		if strings.Contains(ct, "text/html") || strings.Contains(ct, "json") {
			if strings.Contains(rsp.Header.Get("Content-Encoding"), "gzip") {
				g, _ := gzip.NewReader(rsp.Body)
				buf, _ := ioutil.ReadAll(g)
				//buf = bytes.Replace(buf, []byte("陈白痴的博客"), []byte("中国万岁"), -1)
				rsp.Body = ioutil.NopCloser(bytes.NewBuffer(Gzip(buf)))
				fmt.Printf("receive gzip %d data: %#v\n", rsp.StatusCode, string(buf))
			} else {
				buf, _ := ioutil.ReadAll(rsp.Body)
				rsp.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				fmt.Printf("receive %d data: %#v\n", rsp.StatusCode, string(buf))
			}
		}

		if rsp.StatusCode == 301 || rsp.StatusCode == 302 {
			l := rsp.Header.Get("Location")
			if l != "" {
				rsp.Header.Set("Location", strings.Replace(locationplace(l), "https", "http", -1))
			}
		}
		for k, v := range rsp.Header {
			fmt.Println(k, v)
		}
		fmt.Println("**********")
		return nil
	}

	return &httputil.ReverseProxy{Director: director, ModifyResponse: modify}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func Gzip(data []byte) []byte {
	var res bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&res, gzip.DefaultCompression)
	_, err := gz.Write(data)
	if err != nil {
		fmt.Println("zip err" + err.Error())
		return []byte("")
	} else {
		gz.Close()
	}
	return res.Bytes()
}

// 认证的配置
type AwsAuth struct {
	AwsID      string
	AwsKey     string
	AwsRegion  string
	AwsService string
}

func AwsAuthSignature(auth AwsAuth, uri, method, host string, query url.Values, data []byte) (amzdate, authorization_header string) {

	// 基本配置
	access_key := auth.AwsID
	secret_key := auth.AwsKey
	region := auth.AwsRegion
	service := auth.AwsService

	request_parameters := ""

	// 查询字符串排序
	if query != nil {
		temp := []string{}
		for k, _ := range query {
			temp = append(temp, k)
		}
		sort.Strings(temp)
		temp1 := []string{}
		for _, v := range temp {
			temp1 = append(temp1, UriEncode(v, false)+"="+UriEncode(query.Get(v), false))
		}
		request_parameters = strings.Join(temp1, "&")
	}

	// 现在时间
	now := time.Now().UTC()
	amzdate = now.Format("20060102T150405Z")
	datestamp := now.Format("20060102")

	canonical_uri := uri
	canonical_querystring := request_parameters
	canonical_headers := "host:" + host + "\n" + "x-amz-date:" + amzdate + "\n"
	signed_headers := "host;x-amz-date"

	payload_hash := hex.EncodeToString(getSha256Code(""))
	if data != nil {
		payload_hash = hex.EncodeToString(getSha256Code(string(data)))
	}

	canonical_request := method + "\n" + canonical_uri + "\n" + canonical_querystring + "\n" + canonical_headers + "\n" + signed_headers + "\n" + payload_hash
	//fmt.Printf("%q\n", canonical_request)
	algorithm := "AWS4-HMAC-SHA256"

	credential_scope := datestamp + "/" + region + "/" + service + "/" + "aws4_request"
	string_to_sign := algorithm + "\n" + amzdate + "\n" + credential_scope + "\n" + hex.EncodeToString(getSha256Code(canonical_request))
	signing_key := GetSignatureKey(secret_key, datestamp, region, service)
	signature := hex.EncodeToString(Sign(signing_key, []byte(string_to_sign)))
	authorization_header = algorithm + " " + "Credential=" + access_key + "/" + credential_scope + ", " + "SignedHeaders=" + signed_headers + ", " + "Signature=" + signature

	return amzdate, authorization_header
}

func Sign(key, msg []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))
	return h.Sum(nil)
}

func GetSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kDate := Sign([]byte("AWS4"+key), []byte(dateStamp))
	kRegion := Sign(kDate, []byte(regionName))
	kService := Sign(kRegion, []byte(serviceName))
	kSigning := Sign(kService, []byte("aws4_request"))
	return kSigning
}

func getSha256Code(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func UriEncode(src string, encodeSlash bool) string {
	// application/x-www-form-urlencoded will have +
	back := url.QueryEscape(src)
	// all change but + must replace, RFC3986
	temp := strings.Replace(back, "+", "%20", -1)

	// uri / must be keep
	if encodeSlash {
		return strings.Replace(temp, "%2F", "/", -1)
	}
	return temp
}
