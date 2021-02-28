package iot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"sort"
	"strings"
)

func SignMessage(request *resty.Request, sk, ak string) string {
	headerKeys := getSortedHeaders(request)

	headers := strings.ToLower(strings.Join(headerKeys, ";"))
	signature := sign(request, sk)

	return "SDK-HMAC-SHA256 " + "Access=" + ak + ", " + "SignedHeaders=" + headers + "," +
		"Signature=" + signature

}

func sign(request *resty.Request, sk string) string {
	// 计算签名
	h := hmac.New(sha256.New, []byte(sk))
	if _, err := h.Write([]byte(stringToSign(request))); err != nil {
		fmt.Println("error")
	}

	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

func stringToSign(request *resty.Request) string {
	stringToSigns := make([]string, 0)
	stringToSigns = append(stringToSigns, "SDK-HMAC-SHA256")
	stringToSigns = append(stringToSigns, request.Header.Get("X-Sdk-Date"))
	stringToSigns = append(stringToSigns, constructCanonicalRequest(request))

	return strings.Join(stringToSigns, "\n")
}

// 构造规范请求
func constructCanonicalRequest(request *resty.Request) string {
	crParts := make([]string, 0)

	// 构造HTTP请求方法
	httpRequestMethod := strings.ToUpper(request.Method)
	crParts = append(crParts, httpRequestMethod)

	// 添加规范URI参数
	uri, _ := url.Parse(request.URL)
	path := uri.Path
	if len(path) == 0 {
		path = "/"
	}
	if path[len(path)-1] != '/' {
		path += "/"
	}
	canonicalURI := path
	crParts = append(crParts, canonicalURI)

	// 添加规范查询字符串（CanonicalQueryString），以换行符结束。
	queryKeys := make([]string, 0, len(uri.Query()))
	for key := range uri.Query() {
		queryKeys = append(queryKeys, key)
	}
	sort.Strings(queryKeys)

	queryKeysAndValues := make([]string, len(queryKeys))
	for i, key := range queryKeys {
		k := strings.Replace(url.QueryEscape(key), "+", "%20", -1)
		v := strings.Replace(url.QueryEscape(uri.Query().Get(key)), "+", "%20", -1)
		queryKeysAndValues[i] = k + "=" + v
	}

	query := strings.Join(queryKeysAndValues, "&")
	crParts = append(crParts, query)

	// 4 添加规范消息头（CanonicalHeaders），以换行符结束
	headerKeys := getSortedHeaders(request)

	headerKeysAndValues := make([]string, len(headerKeys))
	for i, key := range headerKeys {
		v := request.Header.Get(key)
		k := strings.ToLower(key)
		headerKeysAndValues[i] = k + ":" + v
	}

	crParts = append(crParts, headerKeysAndValues...)
	crParts = append(crParts, "")
	crParts = append(crParts, strings.ToLower(strings.Join(headerKeys, ";")))

	// 使用SHA 256哈希函数以基于HTTP或HTTPS请求正文中的body体（RequestPayload），创建哈希值。
	requestBody, err := getBytes(request.Body)
	if err != nil {
		fmt.Println("get body bytes failed")
		requestBody = []byte("")
	}
	encodedRequestBody := hashRequestBody(requestBody)
	crParts = append(crParts, encodedRequestBody)
	//crParts = append(crParts,"\n")

	stringToHash := strings.Join(crParts, "\n")
	return hashRequestBody([]byte(stringToHash))
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func hashRequestBody(requestBody []byte) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(requestBody))
	if err != nil {
		fmt.Println("error")
	}

	return strings.ToLower(hex.EncodeToString(hash.Sum(nil)))
}

func getSortedHeaders(request *resty.Request) []string {
	headerKeys := make([]string, 0)
	for key := range request.Header {
		headerKeys = append(headerKeys, key)

	}
	sort.Strings(headerKeys)
	return headerKeys
}
