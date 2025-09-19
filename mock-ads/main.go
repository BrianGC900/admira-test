package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		response := map[string]interface{}{
			"external": map[string]interface{}{
				"ads": map[string]interface{}{
					"performance": []map[string]interface{}{
						{
							"date":         "2025-08-01",
							"campaign_id":  "C-1001",
							"channel":      "google_ads",
							"clicks":       1200,
							"impressions":  45000,
							"cost":         350.75,
							"utm_campaign": "back_to_school",
							"utm_source":   "google",
							"utm_medium":   "cpc",
						},
					},
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Mock ADS server running on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
