package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		response := map[string]interface{}{
			"external": map[string]interface{}{
				"crm": map[string]interface{}{
					"opportunities": []map[string]interface{}{
						{
							"opportunity_id": "O-9001",
							"contact_email":  "ana@example.com",
							"stage":          "closed_won",
							"amount":         5000.0,
							"created_at":     time.Now().Format(time.RFC3339),
							"utm_campaign":   "back_to_school",
							"utm_source":     "google",
							"utm_medium":     "cpc",
						},
					},
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Mock CRM server running on :3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
