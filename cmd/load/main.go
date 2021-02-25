package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bbokorney/budget-api-http/pkg/models"
)

func main() {
	fileName := os.Args[1]
	if fileName == "" {
		fmt.Println("Must provide filename")
		os.Exit(1)
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)

	client := http.Client{}
	for {
		tokens, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading line")
			log.Fatal(err)
		}
		if len(tokens) < 4 {
			fmt.Printf("Skipping line, not enough tokens: %s\n", tokens)
			continue
		}
		if tokens[0] == "" {
			fmt.Printf("Skipping line, date is blank: %s\n", tokens)
			continue
		}
		fmt.Println(tokens)

		const longForm = "1/2/2006"
		t, err := time.Parse(longForm, tokens[0])
		if err != nil {
			fmt.Println("Error parsing date")
			log.Fatal(err)
		}
		fmt.Println(t)
		cleanedAmount := strings.Replace(strings.Replace(tokens[1], "$", "", -1), ",", "", -1)
		fmt.Println(cleanedAmount)
		amount, err := strconv.ParseFloat(cleanedAmount, 64)
		if err != nil {
			fmt.Println("Error parsing amount")
			log.Fatal(err)
		}

		transaction := models.Transaction{
			Date:     t,
			Amount:   float32(amount),
			Category: tokens[2],
			Vendor:   tokens[3],
		}
		fmt.Printf("%+v\n", transaction)
		body, err := json.Marshal(&transaction)
		if err != nil {
			fmt.Println("Error marshalling body")
			log.Fatal(err)
		}
		fmt.Println(string(body))
		req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8000/api/v1/transactions", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error creating request")
			log.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Auth-Token", "abc123")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request")
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusAccepted {
			log.Fatalf("Unexpected response code %s", resp.Status)
		}
	}

}
