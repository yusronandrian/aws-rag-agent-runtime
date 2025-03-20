package main

import (
	"aws-rag-agent-runtime/controller"
	"aws-rag-agent-runtime/model"
	"log"
	"net/http"
)

func main() {
	bedrockAgent := model.NewBedrock()
	http.HandleFunc("/send-message", controller.ProcessLLMModel(bedrockAgent))
	log.Println("Server started, listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
