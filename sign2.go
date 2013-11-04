package awsauth

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

func prepareRequestV2(req *http.Request) *http.Request {
	keyID := ""

	if Keys != nil {
		keyID = Keys.AccessKeyID
	}

	values := url.Values{}
	values.Set("AWSAccessKeyId", keyID)
	values.Set("SignatureVersion", "2")
	values.Set("SignatureMethod", "HmacSHA256")
	values.Set("Timestamp", timestampV2())

	augmentRequestQuery(req, values)

	if req.URL.Path == "" {
		req.URL.Path += "/"
	}

	return req
}

func stringToSignV2(req *http.Request) string {
	str := req.Method + "\n"
	str += strings.ToLower(req.URL.Host) + "\n"
	str += req.URL.Path + "\n"
	str += canonicalQueryString(req)
	return str
}

func signatureVersion2(strToSign string) string {
	hashed := hmacSHA256([]byte(Keys.SecretAccessKey), strToSign)
	return base64.StdEncoding.EncodeToString(hashed)
}

func augmentRequestQuery(req *http.Request, values url.Values) *http.Request {
	for key, arr := range req.URL.Query() {
		for _, val := range arr {
			values.Set(key, val)
		}
	}

	req.URL.RawQuery = values.Encode()

	return req
}

func canonicalQueryString(req *http.Request) string {
	return req.URL.RawQuery
}

func timestampV2() string {
	return now().Format(timeFormatV2)
}

const timeFormatV2 = "2006-01-02T15:04:05"
