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
