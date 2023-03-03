package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/guotie/config"
	"github.com/swgloomy/gutil/glog"
	"io"
	"net/http"
	"net/url"
)

type chatGptStruct struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

type chatGptResponseStruct struct {
	Id      string                 `json:"id"`
	Object  string                 `json:"object"`
	Model   string                 `json:"model"`
	Choices []choicesContentStruct `json:"choices"`
}

type choicesContentStruct struct {
	Text string `json:"text"`
}

func AcquireContent(message string) string {
	content, err := completions(message)
	if err != nil {
		glog.Error("package:chatgpt func:AcquireContent completions run err! err: %+v \n", err)
		return "啊哦,我出错了.怎么办,好苦恼"
	}
	return content
}

func completions(msg string) (string, error) {
	var payload chatGptStruct
	payload.Model = "text-davinci-003"
	payload.Prompt = msg
	payload.Temperature = 0.7
	payload.MaxTokens = 256
	payload.TopP = 1
	payload.FrequencyPenalty = 0
	payload.PresencePenalty = 0
	bs, err := json.Marshal(payload)
	if err != nil {
		glog.Error("package:chatgpt func:completions Marshal run err! model: %+v err: %+v \n", payload, err)
		return "", err
	}
	uri, err := url.Parse("socket5://127.0.0.1:10808")
	if err != nil {
		glog.Error("package:chatgpt func:completions url parse run err! model: %+v err: %+v \n", payload, err)
		return "", err
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", bytes.NewReader(bs))
	if err != nil {
		glog.Error("package:chatgpt func:completions NewRequest run err! err: %+v \n", payload, err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.GetString("chatGptSecret"))
	resp, err := client.Do(req)
	if err != nil {
		glog.Error("package:chatgpt func:completions client do run err! model: %+v err: %+v \n", payload, err)
		return "", err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			glog.Error("package:chatgpt func:completions body close err! model: %+v err: %+v \n", payload, err)
		}
	}()

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		glog.Error("package:chatgpt func:completions ReadAll body do run err! model: %+v err: %+v \n", payload, err)
		return "", err
	}
	var (
		responseModel   chatGptResponseStruct
		responseContent string
	)
	err = json.Unmarshal(bodyByte, &responseModel)
	if err != nil {
		glog.Error("package:chatgpt func:completions Unmarshal run err! response: %s err: %+v \n", string(bodyByte), err)
		return "", err
	}
	for _, choice := range responseModel.Choices {
		responseContent += fmt.Sprintf("%s\r\n", choice.Text)
	}
	return responseContent, nil
}
