package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strings"
)

var playerName string

func register() {
	var name string
	fmt.Print("Enter your name to register: ")
	fmt.Scanln(&name)

	data := map[string]string{"name": name}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error registering:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Registration successful!")
		playerName = name
		showMainMenu()
	} else {
		fmt.Println("Failed to register:", resp.Status)
	}
}

func login() {
	var name string
	fmt.Print("Enter your name to login: ")
	fmt.Scanln(&name)

	data := map[string]string{"name": name}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error logging in:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Login successful!")
		playerName = name
		showMainMenu()
	} else {
		fmt.Println("Failed to login:", resp.Status)
	}
}

func createWorld() {
	resp, err := http.Post("http://localhost:8080/createWorld", "application/json", nil)
	if err != nil {
		fmt.Println("Error creating world:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response struct {
			Message      string         `json:"message"`
			WorldID      int            `json:"world_id"`
			ActiveWorlds map[int]string `json:"active_worlds"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}
		fmt.Println(response.Message)
		fmt.Println("Active worlds:")
		for id, desc := range response.ActiveWorlds {
			fmt.Printf("World ID: %d - %s\n", id, desc)
		}
	} else {
		fmt.Println("Failed to create world:", resp.Status)
	}
}

func joinWorld() {
	resp, err := http.Get("http://localhost:8080/getActiveWorlds")
	if err != nil {
		fmt.Println("Error fetching active worlds:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var activeWorlds map[int]string
		err = json.NewDecoder(resp.Body).Decode(&activeWorlds)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}

		fmt.Println("Active worlds:")
		for id, desc := range activeWorlds {
			fmt.Printf("World ID: %d - %s\n", id, desc)
		}

		var worldID int
		fmt.Print("Enter the World ID to join: ")
		fmt.Scan(&worldID)

		var mode string
		fmt.Print("Choose movement mode (auto/manual): ")
		fmt.Scan(&mode)

		data := map[string]interface{}{
			"action":      "join_world",
			"player_name": playerName,
			"world_id":    worldID,
			"mode":        mode,
		}

		if mode == "auto" {
			var delay int
			fmt.Print("Enter auto move delay (ms): ")
			fmt.Scan(&delay)
			data["auto_move_delay"] = delay
		} else if mode == "manual" {
			var direction string
			fmt.Print("Enter direction (up/down/left/right): ")
			fmt.Scan(&direction)
			data["direction"] = direction
		}

		jsonData, _ := json.Marshal(data)

		c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		err = c.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("write:", err)
			return
		}

		messageChan := make(chan []byte)

		go func() {
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					close(messageChan)
					return
				}
				messageChan <- message
			}
		}()

		fmt.Println("Listening for updates...")
		for message := range messageChan {
			fmt.Printf("Message from server: %s\n", message)

			if strings.Contains(string(message), "Do you want to catch it? (yes/no):") {
				var answer string
				fmt.Print("Enter your choice (yes/no): ")
				fmt.Scan(&answer)

				response := map[string]string{
					"response": answer,
				}

				responseData, err := json.Marshal(response)
				if err != nil {
					log.Println("error marshalling response:", err)
					continue
				}

				err = c.WriteMessage(websocket.TextMessage, responseData)
				if err != nil {
					log.Println("write:", err)
					continue
				}
			}
		}
	} else {
		fmt.Println("Failed to fetch active worlds:", resp.Status)
	}
}

func showMainMenu() {
	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Create World")
		fmt.Println("2. Join World")
		fmt.Println("3. Exit")
		var choice int
		fmt.Scan(&choice)
		switch choice {
		case 1:
			createWorld()
		case 2:
			joinWorld()
		case 3:
			os.Exit(0)
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func main() {
	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		var choice int
		fmt.Scan(&choice)
		switch choice {
		case 1:
			register()
		case 2:
			login()
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
