/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
)

func mockCardAction() []byte {
	var token = "v"
	value := map[string]interface{}{}
	value["value"] = "sdfsfd"
	value["tag"] = "button"

	cardAction := &larkcard.CardAction{
		OpenID:        "ou_sdfimx9948345",
		UserID:        "eu_sd923r0sdf5",
		OpenMessageID: "om_abcdefg1234567890",
		TenantKey:     "d32004232",
		Token:         token,
		Action: &struct {
			Value    map[string]interface{} `json:"value"`
			Tag      string                 `json:"tag"`
			Option   string                 `json:"option"`
			Timezone string                 `json:"timezone"`
		}{
			Value: value,
			Tag:   "button",
		},
	}

	cardActionBody := &larkcard.CardActionBody{
		CardAction: cardAction,
		Challenge:  "121212",
		Type:       "url_verification1",
	}

	body, _ := json.Marshal(cardActionBody)
	return body
}
func main() {

	body1 := "{\"open_id\":\"ou_cdaf8ff9dcf6723db8f09f10b4b53c0f\",\"open_message_id\":\"om_697921c87902fa9999cf9b8520ea5eb2\",\"open_chat_id\":\"oc_79dd04bf25e01d267faadb93b11caac7\",\"tenant_key\":\"104740a13d02575e\",\"token\":\"c-bb72b8787895b65b7fd5ec2ce18aff2bf8cad2af\",\"action\":{\"value\":{\"cardData\":\"{\\\"actionApproveId\\\":5483,\\\"approvalType\\\":0,\\\"approveDetailUrl\\\":\\\"https://wxpublic-test2.t3go.cn/t3-h5-company/#/approveBridge?redirectUrl=https%3A%2F%2Fwxpublic-test2.t3go.cn%2Fcompany-approve-h5%2F%23%2FapproveDetailNew%3FapproveId%3D5483\\\",\\\"approverId\\\":\\\"af075306825d49358ba4662e7c6d9488\\\",\\\"approverIds\\\":[\\\"af075306825d49358ba4662e7c6d9488\\\"],\\\"employeeId\\\":\\\"3a00f21601e748da8a6fa97e393557c3\\\",\\\"employeeName\\\":\\\"王树轩\\\",\\\"orgId\\\":\\\"75e1419b996d4e6da70056e8e1a42ee8\\\",\\\"reason\\\":\\\"Jhh\\\",\\\"routeDetail\\\":\\\"A(花园城巨石马群店)--B座\\\",\\\"sceneShowName\\\":\\\"审批\\\",\\\"thirdAppType\\\":1,\\\"thirdEmployeeId\\\":\\\"ou_6a73fb461db56e87fdb7c922f353dddb\\\",\\\"useDate\\\":\\\"2022/11/15 00:00-2022/11/16 00:00\\\"}\",\"handleType\":\"1\"},\"tag\":\"button\"}}"
	sourceSign := larkcard.Signature("Monday, 14-Nov-22 15:53:37 CST", "39ef4f6f-5400-41e4-b18d-0d13ab6d3646", "UtVPQSv5lbtzjvaulJjvPg8RAOmXTbGF", body1)
	fmt.Println(sourceSign)
	var b strings.Builder
	b.WriteString("Monday, 14-Nov-22 15:53:37 CST")
	b.WriteString("39ef4f6f-5400-41e4-b18d-0d13ab6d3646")
	b.WriteString("UtVPQSv5lbtzjvaulJjvPg8RAOmXTbGF")
	b.WriteString(body1) //body指整个请求体，不要在反序列化后再计算
	bt := []byte(b.String())
	h := sha1.New()
	h.Write(bt)
	bs := h.Sum(nil)
	sig := fmt.Sprintf("%x", bs)
	fmt.Println(sig)
	// check if request headers['X-Lark-Signature'] equals to signature

	//mock body
	body := mockCardAction()

	// 创建http req
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:7777/webhook/card", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 计算签名
	var timestamp = "timestamp"
	var nonce = "nonce"
	var token = "v"

	//var b = "{\"open_id\":\"ou_d840b2e2be16b3e0091bc0c79220e1fa\",\"user_id\":\"16fd348g\",\"open_message_id\":\"om_dce5707d696ee4952ebedaf1ee762ed2\",\"tenant_key\":\"736588c9260f175d\",\"token\":\"v\",\"action\":{\"value\":{\"key\":\"value\"},\"tag\":\"button\"}}"
	sourceSign = larkcard.Signature(timestamp, nonce, token, string(body))
	//fmt.Println(sourceSign)
	// 添加header
	req.Header.Set(larkevent.EventRequestTimestamp, timestamp)
	req.Header.Set(larkevent.EventRequestNonce, nonce)
	req.Header.Set(larkevent.EventSignature, sourceSign)
	req.Header.Set("X-Tt-Logid", "logid111111111111111")

	// 模拟推送卡片消息
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 结果处理
	fmt.Println(resp.StatusCode)
	fmt.Println(larkcore.Prettify(resp.Header))
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
