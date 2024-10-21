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

var testA = make([]string, 13) // Aloca as primeiras 13 letras do alfabeto

func init() {
	// Inicializa o array com as primeiras 13 letras do alfabeto
	for i := range testA {
		testA[i] = string('a' + i)
	}
}

// Função para encaminhar requisição à Communication API
func forwardToCommunicationAPI(option int, letter string) (string, error) {
	logRequest := LogMessage{Option: option, Letter: letter}
	jsonData, err := json.Marshal(logRequest)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Faz a requisição à Communication API
	resp, err := http.Post("http://communication_api:8083/forward", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request to Communication API: %v", err)
	}
	defer resp.Body.Close()

	// Ler a resposta da Communication API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response from Communication API: %v", err)
	}

	// Retorna a resposta como string
	return string(body), nil
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	var logMsg LogMessage
	err := json.NewDecoder(r.Body).Decode(&logMsg)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Verifica se a opção é para uma letra que está no servidor A
	if logMsg.Option < 0 || logMsg.Option >= 26 {
		http.Error(w, "Invalid option", http.StatusBadRequest)
		return
	}

	// Se a letra está entre 0 e 12, lidamos localmente
	if logMsg.Option < 13 {
		// Se não houver uma nova letra, apenas retorne o valor atual
		if logMsg.Letter == "" {
			response := map[string]string{
				"message": testA[logMsg.Option],
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Atualizar o valor de testA[option]
		testA[logMsg.Option] = logMsg.Letter
		response := map[string]string{
			"message": fmt.Sprintf("Letter at position %d updated to '%s'", logMsg.Option, logMsg.Letter),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Se a letra está entre 13 e 25, encaminhamos para a Communication API
	message, err := forwardToCommunicationAPI(logMsg.Option, logMsg.Letter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder ao cliente com a resposta da Communication API
	fmt.Fprintf(w, "%s\n", message)
}

func main() {
	http.HandleFunc("/log", logHandler)
	fmt.Println("Server A is listening on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
