package cgv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
)

const CgvAddr = "http://www.cgv.com.cn"

type CommonResponse struct {
	SVC_ERR_MSG_TEXT string
	RS_MSG           string
	RS_CD            string
	ErrorCode        int
}

func request(method, path string, p url.Values, reqInfo interface{}, respInfo interface{}) error {
	var body io.Reader
	if reqInfo != nil {
		dat, err := json.Marshal(reqInfo)
		if err != nil {
			return fmt.Errorf("编码请求参数失败: ", err)
		}

		body = bytes.NewReader(dat)
	}

	req, err := http.NewRequest(method, CgvAddr+path+"?"+p.Encode(), body)
	if err != nil {
		return fmt.Errorf("创建请求失败: ", err)
	}

	if reqInfo != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("执行请求失败: ", err)
	}
	defer resp.Body.Close()

	mediaType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if mediaType != "application/json" {
		dat, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(resp.Status + "\n" + string(dat))
	}

	buf, err := ioutil.ReadAll(resp.Body)

	var info CommonResponse
	err = json.Unmarshal(buf, &info)
	if err != nil {
		return fmt.Errorf("解析通用响应失败: %s", err)
	}

	//TODO: 判断公共响应中的错误

	err = json.Unmarshal(buf, respInfo)
	if err != nil {
		return fmt.Errorf("解析业务响应失败: %s", err)
	}

	return nil
}
