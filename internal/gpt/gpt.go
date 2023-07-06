package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiEndpoint = "https://api.openai.com/v1/engines/davinci/completions"
)

type GPTRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type GPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type GPTRequester struct {
	key string
}

func NewGPTRequester(key string) *GPTRequester {
	return &GPTRequester{key}
}

func (g *GPTRequester) CheckCorrectness(sentence string) (string, error) {
	return g.doGPTRequest(fmt.Sprintf(`Check the correctness of sentence: "%s"`, sentence))
}

func (g *GPTRequester) FindContextSentence(context string) (string, error) {
	return g.doGPTRequest(fmt.Sprintf(`Create an example sentence containing words: "%s"`, context))
}

func (g *GPTRequester) doGPTRequest(question string) (string, error) {
	// Prepare the request payload
	gptRequest := GPTRequest{
		Prompt:      question,
		MaxTokens:   300,
		Temperature: 0.7,
	}
	payload, err := json.Marshal(gptRequest)
	if err != nil {
		return "", err
	}

	// Send the request to the API
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.key))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse the response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var gptResponse GPTResponse

	err = json.Unmarshal(respBody, &gptResponse)
	if err != nil {
		return "", err
	}

	// Extract the generated sentence from the response
	if len(gptResponse.Choices) > 0 {
		return gptResponse.Choices[0].Text, nil
	}
	return "", nil
}
