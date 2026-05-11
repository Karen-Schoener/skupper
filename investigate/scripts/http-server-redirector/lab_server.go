/* IBM Confidential                 */
/* PID 5900-BXU                     */
/* © Copyright IBM Corp. 2026       */

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// printRequestMetadata dumps the method, URL, and all headers to the console.
func printRequestMetadata(serverName string, r *http.Request) {
	fmt.Printf("\n--- %s Incoming Request ---\n", serverName)
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL:    %s\n", r.URL.String())
	fmt.Printf("Headers:\n")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}
}

func main() {
	// --- SERVER 1: The Redirector (Port 8081) ---
	go func() {
		mux1 := http.NewServeMux()
		mux1.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			printRequestMetadata("SERVER 8081", r)
			
			// Determine status code from query param, default to 307
			statusParam := r.URL.Query().Get("status")
			statusCode := http.StatusTemporaryRedirect // 307

			switch statusParam {
			case "301":
				statusCode = http.StatusMovedPermanently
			case "302":
				statusCode = http.StatusFound
			case "303":
				statusCode = http.StatusSeeOther
			case "308":
				statusCode = http.StatusPermanentRedirect
			}

			// Capture and print the body at the first hop
			body, _ := io.ReadAll(r.Body)
			if len(body) > 0 {
				fmt.Printf("Body:   %s\n", string(body))
			} else {
				fmt.Printf("Body:   [Empty]\n")
			}

			fmt.Printf("Action: Sending %d Redirect to :8082/final\n", statusCode)
			
			w.Header().Set("Location", "http://localhost:8082/final")
			w.WriteHeader(statusCode)
		})

		log.Println("Starting Server 1 (Redirector) on :8081...")
		if err := http.ListenAndServe(":8081", mux1); err != nil {
			log.Fatalf("Server 1 failed: %v", err)
		}
	}()

	// --- SERVER 2: The Final Destination (Port 8082) ---
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/final", func(w http.ResponseWriter, r *http.Request) {
		printRequestMetadata("SERVER 8082", r)

		// Explicitly check for our sensitive test headers
		fmt.Printf("Auth Header:  %s\n", r.Header.Get("Authorization"))
		fmt.Printf("X-Auth-Token: %s\n", r.Header.Get("X-Auth-Token"))

		// Read and verify if the body survived the redirect
		body, _ := io.ReadAll(r.Body)
		if len(body) > 0 {
			fmt.Printf("Body:         %s\n", string(body))
		} else {
			fmt.Printf("Body:         [EMPTY - PAYLOAD LOST!]\n")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Target Reached\n"))
		fmt.Println("Action: Responded with 200 OK")
	})

	log.Println("Starting Server 2 (Target) on :8082...")
	log.Println("Lab ready. Use ?status=301, 307, or 308 to test different redirect behaviors.")
	
	// This blocks the main thread
	if err := http.ListenAndServe(":8082", mux2); err != nil {
		log.Fatalf("Server 2 failed: %v", err)
	}
}
