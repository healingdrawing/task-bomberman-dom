package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
)

// # jsonResponse marshals and forwards json response writing to http.ResponseWriter
//
// @params {w http.ResponseWriter, statusCode int, data any}
// @sideEffect {jsonResponse -> w}
func jsonResponse(w http.ResponseWriter, statusCode int, data any) {
	// only status "ok"(200) allows to show "message" content in browser console. With other statuses browser displays only status code. Maybe it is golang http package "feature"

	if message, ok := data.(string); ok { // if data type is string
		jsonResponseObj, _ := json.Marshal(map[string]string{
			"message": http.StatusText(statusCode) + ": " + message,
		})
		w.WriteHeader(statusCode)

		_, err := w.Write(jsonResponseObj)
		if err != nil {
			log.Println(err.Error())
		}

	} else { // if len(jsonResponseObj) == 0 { // if unhandled by above custom conversion
		w.WriteHeader(statusCode)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

}

// # recovery is a utility function to recover from panic and send a json err response over http
//
// @sideEffect {log, debug}
//
// - for further debugging uncomment {print stack trace}
func recovery(w http.ResponseWriter) {
	if r := recover(); r != nil {
		fmt.Println("=====================================")
		stackTrace := debug.Stack()
		lines := strings.Split(string(stackTrace), "\n")
		relevantPanicLines := []string{}
		for _, line := range lines {
			if strings.Contains(line, "backend/") {
				relevantPanicLines = append(relevantPanicLines, line)
			}
		}
		if len(relevantPanicLines) > 1 {
			for i, line := range relevantPanicLines {
				if strings.Contains(line, "utils.go") {
					relevantPanicLines = append(relevantPanicLines[:i], relevantPanicLines[i+1:]...)
				}
			}
		}
		relevantPanicLine := strings.Join(relevantPanicLines, "\n")
		log.Println(relevantPanicLines)
		jsonResponse(w, http.StatusInternalServerError, relevantPanicLine)
		fmt.Println("=====================================")
		// to print the full stack trace
		log.Println(string(stackTrace))
	}
}

// randomNum returns a random number between min and max, both inclusive.
func randomNum(min, max int) int {
	bi := big.NewInt(int64(max + 1 - min))
	bj, err := rand.Int(rand.Reader, bi)
	if err != nil {
		log.Fatal(err)
	}
	return int(bj.Int64()) + min
}

func generateRandomEmojiSequence() string {
	rounds := []string{"🔴", "🟠", "🟡", "🟢", "🔵", "🟣", "🟤", "⚫", "⚪"}
	// Shuffle the rounds using Fisher-Yates algorithm
	for i := len(rounds) - 1; i > 0; i-- {
		bi := big.NewInt(3)
		bj, err := rand.Int(rand.Reader, bi)
		if err != nil {
			log.Fatal(err)
		}
		// convert big.Int to int
		j := int(bj.Int64())
		rounds[i], rounds[j] = rounds[j], rounds[i]
	}

	// Join the shuffled rounds into a single string
	mixedRounds := strings.Join(rounds, " ")

	return mixedRounds
}

func generate_UUID() (UUID string, err error) {
	random, _ := uuid.NewV4()
	UUID = random.String()
	return UUID, nil
}

// connected_clients_number counts the clients number in a sync.Map
func connected_clients_number(clients *sync.Map) int {
	count := 0
	clients.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// choose_first_free_number returns the first free number in a sync.Map of four elements. The clien.NUMBER will be checked for each client in the map. Used inside server , to send to client instructions about color and number of player/corner of starting position
func choose_first_free_number(clients *sync.Map) int {
	used_numbers := map[int]int{}
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		used_numbers[client.NUMBER] = client.NUMBER
		return true
	})

	for i := 1; i < 5; i++ {
		if _, ok := used_numbers[i]; !ok {
			return i
		}
	}
	return 0
}

// not beautiful, every step will happens iteration over all clients, even if it is only x4 clients maximum
func get_all_clients_uuids(clients *sync.Map) []string {
	var uuids []string
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		uuids = append(uuids, client.UUID)
		return true
	})
	return uuids
}

// get_client_by_uuid returns a client by uuid
func get_client_by_uuid(clients *sync.Map, uuid string) *Client {
	var client *Client
	clients.Range(func(key, value interface{}) bool {
		client = value.(*Client)
		if client.UUID == uuid {
			return false
		}
		return true
	})
	return client
}

// get client number by uuid
func get_client_number_by_uuid(clients *sync.Map, uuid string) int {
	return get_client_by_uuid(clients, uuid).NUMBER
}

// get client nickname by uuid
func get_client_nickname_by_uuid(clients *sync.Map, uuid string) string {
	return get_client_by_uuid(clients, uuid).NICKNAME
}
