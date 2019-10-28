package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)
func main() {
	file := os.Getenv("DUNDERGITCALL_FILE")

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("[ERR] %v", err)
		os.Exit(1)
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(b, &data); err != nil {
		log.Printf("[ERR] %v", err)
		os.Exit(1)
	}

	data["timestamp"] = time.Now().Unix()

	b, err = json.Marshal(data)
	if err != nil {
		log.Printf("[ERR] %s", err)
		os.Exit(1)
	}

	if webhook, ok := data["webhook"].(string); ok {
		req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(b))
		if err != nil {
			log.Printf("[ERR] %s", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[ERR] %s", err)
			os.Exit(1)
		}
		resp.Body.Close()
	}

	if err := ioutil.WriteFile(file, b, 0644); err != nil {
		log.Printf("[ERR] %s", err)
		os.Exit(1)
	}
}
