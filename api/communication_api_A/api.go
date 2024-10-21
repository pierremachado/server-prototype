package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type LogMessage struct {
	Option int    `json:"option"`
	Letter string `json:"letter,omitempty"`
}

// Estrutura para representar a resposta do servidor B
type ResponseMessage struct {
	Message string `json:"message"`
}

// Função para encaminhar requisição ao servidor B
func forwardToServerB(option int, letter string) (string, error) {
	logRequest := LogMessage{Option: option, Letter: letter}
	jsonData, err := json.Marshal(logRequest)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Faz a requisição ao servidor B
	resp, err := http.Post("http://server_b:8082/log", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request to server B: %v", err)
	}
	defer resp.Body.Close()

	// Verifica se o status da resposta é 200 OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: server B returned status code %d", resp.StatusCode)
	}

	// Ler a resposta do servidor B
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response from server B: %v", err)
	}

	// Fazer unmarshal do JSON recebido para a estrutura ResponseMessage
	var response ResponseMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response from server B: %v", err)
	}

	// Retorna apenas o valor da mensagem extraída
	return response.Message, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var logMsg LogMessage
	err := json.NewDecoder(r.Body).Decode(&logMsg)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Encaminha a requisição ao servidor B
	message, err := forwardToServerB(logMsg.Option, logMsg.Letter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define o Content-Type como JSON
	w.Header().Set("Content-Type", "application/json")

	// Retorna a resposta em JSON para o servidor A, contendo apenas a mensagem
	response := map[string]string{
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/forward", handleRequest)
	fmt.Println("Communication API is listening on port 8083...")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
