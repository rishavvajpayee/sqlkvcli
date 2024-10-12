package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const serverURL = "http://localhost:8000"

func processCommand(command string, args []string) {
	switch command {
	case "get":

		if len(args) != 1 || len(args) == 0 {
			fmt.Println("Usage: get <key>")
			return
		}
		key := args[0]
		resp, err := http.Get(fmt.Sprintf("%s/kv/get/%s", serverURL, key))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error:", resp.Status)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		var result = map[string]interface{}{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(result[key])
	case "set":
		if len(args) < 2 {
			fmt.Println("Usage: set <key> <value> <exp>")
			return
		}
		exp := ""
		key := args[0]
		value := args[1]
		if len(args) == 3 {
			exp = args[2]
		}
		jsonBody := []byte(nil)
		if exp != "" {
			jsonBody = []byte(`{"key": "` + key + `", "value": ` + value + `, "expires_in": ` + exp + `}`)
		} else {
			jsonBody = []byte(`{"key": "` + key + `", "value": ` + value + `}`)
		}
		request, err := http.NewRequest("POST", fmt.Sprintf("%s/kv/set", serverURL), bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	case "exit":
		fmt.Println("Exiting CLI...")
		os.Exit(0)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the interactive KV Store CLI")
	fmt.Println("Available commands: get <key>, set <key> <value>, exit")
	for {
		fmt.Print("kvcli > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")
		if len(parts) == 0 {
			fmt.Println("Invalid command")
			continue
		}
		command := parts[0]
		args := parts[1:]

		processCommand(command, args)
	}
}
