package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type LogRequest struct {
	Option int    `json:"option"`
	Letter string `json:"letter,omitempty"`
}

func main() {
	for {
		fmt.Println("Enter a number to request or update the server (0-25):")
		var option int
		fmt.Scanln(&option)

		if option < 0 || option > 25 {
			fmt.Println("Invalid option. Please enter a number between 0 and 25.")
			continue
		}

		fmt.Println("Do you want to update the letter at this position? (y/n):")
		var update string
		fmt.Scanln(&update)

		var letter string
		if update == "y" || update == "Y" {
			fmt.Println("Enter the new letter:")
			fmt.Scanln(&letter)
		}

		// Montar o corpo da requisição JSON
		logRequest := LogRequest{Option: option, Letter: letter}
		jsonData, err := json.Marshal(logRequest)
		if err != nil {
			log.Fatalf("Error marshalling JSON: %v", err)
		}

		// Fazer a requisição POST para a API
		resp, err := http.Post("http://client_api:8080", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		// Ler a resposta da API
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		fmt.Println("Response from API:", string(body))

		// Pequeno atraso antes do próximo prompt
		time.Sleep(1 * time.Second)
	}
}
