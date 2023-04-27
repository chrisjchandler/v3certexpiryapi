package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the file path from the request parameters
		filePath := r.URL.Query().Get("file_path")

		// Read the certificate file
		certData, err := ioutil.ReadFile(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Decode the PEM-encoded certificate data to DER-encoded data, if necessary
		block, _ := pem.Decode(certData)
		if block != nil {
			certData = block.Bytes
		}

		// Parse the certificate
		cert, err := x509.ParseCertificate(certData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate the number of days remaining before the certificate expires
		now := time.Now()
		duration := cert.NotAfter.Sub(now)
		daysRemaining := int(duration.Hours() / 24)

		// Return the number of days remaining as a JSON response
		w.Header().Set("Content-Type", "application/json")
		response := map[string]int{"days_remaining": daysRemaining}
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", nil)
}
