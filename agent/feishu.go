package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joycastle/casual-server-lib/log"
)

type Request struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Text string `json:"text"`
}

//发送给后端
func FeiShuSendToServiceRD(msg string) error {
	var req Request
	req.MsgType = "text"
	req.Content.Text = msg

	bodyJson, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("msg:%s, error:%s", msg, err.Error())
	}

	url := "https://open.feishu.cn/open-apis/bot/v2/hook/560c4814-30ef-4ff7-90e6-e2fa6008c6ec"

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")

	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()

	rspBody, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return err
	}

	log.Get("agent").Info("Agent-Feishu", "request:", string(bodyJson), "respone:", string(rspBody))

	return nil
}
