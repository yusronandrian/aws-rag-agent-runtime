package controller

import (
	"aws-rag-agent-runtime/model"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type LLMRequest struct {
	Question string `json:"question"`
}

type LLMResponse struct {
	Response string `json:"response"`
}

func ProcessLLMModel(bedrockAgent *model.BedrockAgent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method should be POST!", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Parse the request body as JSON
		var req LLMRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		question := strings.TrimSpace(req.Question)

		// Check if the question field is empty.
		if question == "" {
			http.Error(w, "The 'question' field is required", http.StatusBadRequest)
			return
		}

		// Pass the question to Knowledge Base and return back the response.
		response := bedrockAgent.RetrieveResponseFromKnowledgeBase(question)
		var llmResponse LLMResponse
		llmResponse.Response = response

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(llmResponse)
	}
}
