package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaRequest struct {
	Model     string          `json:"model"`
	Stream    bool            `json:"stream"`
	Messages  []OllamaMessage `json:"messages"`
	KeepAlive int             `json:"keep_alive"`
}

type OllamaResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Message            OllamaMessage `json:"message"`
	DoneReason         string        `json:"done_reason"`
	Done               bool          `json:"done"`
	TotalDuration      int64         `json:"total_duration"`
	LoadDuration       int64         `json:"load_duration"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration int64         `json:"prompt_eval_duration"`
	EvalCount          int           `json:"eval_count"`
	EvalDuration       int64         `json:"eval_duration"`
}

func getAIResponse(messages []OllamaMessage) string {
	url := "http://localhost:11434/api/chat"

	requestBody, err := json.Marshal(OllamaRequest{
		Model:     "mannix/llama3.1-8b-lexifsfsdfsdf:tools-q8_0",
		Stream:    false,
		Messages:  messages,
		KeepAlive: 10000,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	fmt.Println("Response:", resp)
	if err != nil {
		return "<h1>Error: " + err.Error() + "</h1>"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	var errorResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
		return "<h1>Error: " + errorResp.Error + "</h1>"
	}
	var ollamaResp OllamaResponse
	err = json.Unmarshal(body, &ollamaResp)
	if err != nil {
		return ""
	}
	return ollamaResp.Message.Content
}
