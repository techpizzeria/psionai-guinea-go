// Command guinea-completion is a tiny HTTP service that answers prompts by
// forwarding them to the OpenAI completions API.
//
// POST /complete with a JSON body {"prompt": "..."} returns {"text": "..."}.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/techpizzeria/psionai-guinea-go/openai"
)

// completeRequest is the body accepted by the /complete endpoint.
type completeRequest struct {
	Prompt string `json:"prompt"`
}

// completeResponse is the body returned by the /complete endpoint.
type completeResponse struct {
	Text string `json:"text"`
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}
	client := openai.NewClient(apiKey)

	http.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req completeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Prompt) == "" {
			http.Error(w, "prompt is required", http.StatusBadRequest)
			return
		}

		text, err := client.Complete(r.Context(), req.Prompt)
		if err != nil {
			log.Printf("completion failed: %v", err)
			http.Error(w, "upstream completion failed", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(completeResponse{Text: text}); err != nil {
			log.Printf("write response: %v", err)
		}
	})

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
