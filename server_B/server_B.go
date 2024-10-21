package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LogMessage struct {
	Option int    `json:"option"`
	Letter string `json:"letter,omitempty"`
}

var testB = make([]string, 13) // Aloca as últimas 13 letras do alfabeto

func init() {
	// Inicializa o array com as últimas 13 letras do alfabeto
	for i := range testB {
		testB[i] = string('n' + i) // 'n' é a 14ª letra do alfabeto
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received log message")

	var logMsg LogMessage
	err := json.NewDecoder(r.Body).Decode(&logMsg)
	if err != nil {
		// Retorna um erro em formato JSON
		fmt.Println("E1")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Verifica se a opção está dentro do intervalo de 13 a 25 (que é o testB)
	if logMsg.Option < 13 || logMsg.Option >= 26 {
		// Retorna um erro em formato JSON
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("E2")
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid option for Server B"})
		return
	}

	// Convertendo para o índice do array testB
	index := logMsg.Option - 13

	// Se não houver uma nova letra, apenas retorne o valor atual de testB
	if logMsg.Letter == "" {
		response := map[string]string{
			"message": testB[index],
		}
		fmt.Println("E3")
		w.Header().Set("Content-Type", "application/json") // Define o tipo de resposta como JSON
		json.NewEncoder(w).Encode(response)
		return
	}

	// Atualizar o valor de testB[index] com a nova letra
	testB[index] = logMsg.Letter
	response := map[string]string{
		"message": fmt.Sprintf("Letter at position %d updated to '%s'", logMsg.Option, logMsg.Letter),
	}
	w.Header().Set("Content-Type", "application/json") // Define o tipo de resposta como JSON
	json.NewEncoder(w).Encode(response)

	fmt.Println(testB)
}

func main() {
	http.HandleFunc("/log", logHandler)
	fmt.Println("Server B is listening on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
