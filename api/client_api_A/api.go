package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type LogRequest struct {
	Option int    `json:"option"`
	Letter string `json:"letter,omitempty"`
}

func sendToLogServer(option int, letter string) (string, error) {
	// Enviar a opção e a letra ao servidor de log
	logRequest := LogRequest{Option: option, Letter: letter}
	jsonData, err := json.Marshal(logRequest)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Faz a requisição ao servidor
	resp, err := http.Post("http://server_a:8081/log", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Ler a resposta do servidor
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Extrair a mensagem da resposta JSON
	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	return result["message"], nil
}

func helloMomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Ler o corpo da requisição do cliente
	var logRequest LogRequest
	err := json.NewDecoder(r.Body).Decode(&logRequest)
	if err != nil || logRequest.Option < 0 || logRequest.Option > 25 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Enviar a opção (e letra, se fornecida) ao servidor
	message, err := sendToLogServer(logRequest.Option, logRequest.Letter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder ao cliente com a mensagem recebida
	fmt.Fprintf(w, "API Server: %s\n", message)
}

func main() {
	http.HandleFunc("/", helloMomHandler)
	fmt.Println("API Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
