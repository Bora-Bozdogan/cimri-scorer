package handlers

import (
	"bytes"
	"encoding/json"
	"scorer/internal/models"
	"net/http"
	"testing"
)

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }

func TestHandleUpdate(t *testing.T) {
	MTest.TestHandleUpdate(t)
}

func (m *MainTest) TestHandleUpdate(t *testing.T) {
	app := m.app

	tests := []struct {
		name     string
		request  *models.Request
		rawBody  []byte
		apiKey   string
		expected int
		contains string
	}{
		{
			name: "Minor price drop, existing product",
			request: &models.Request{
				ApiKey:             strPtr("amazon-key"),
				ProductName:        strPtr("IPhone 16"),
				ProductDescription: strPtr("Latest Apple flagship phone"),
				ProductImage:       strPtr("https://example.com/iphone16.jpg"),
				StoreName:          strPtr("Amazon"),
				Price:              intPtr(1475),
				Stock:              intPtr(50),
				PopularityScore:    intPtr(5),
				UrgencyScore:       intPtr(5),
			},
			apiKey:   "amazon-key",
			expected: http.StatusOK, // 10*30 + 5*25 + 5*20 + 5*15 = 300 + 125 + 100 + 75 = 600
			contains: "600 sent to que",
		},
		{
			name: "Stock out, new merchant",
			request: &models.Request{
				ApiKey:             strPtr("trendyol-key"),
				ProductName:        strPtr("IPhone 16"),
				ProductDescription: strPtr("Latest Apple flagship phone"),
				ProductImage:       strPtr("https://example.com/iphone16.jpg"),
				StoreName:          strPtr("Trendyol"),
				Price:              intPtr(1500),
				Stock:              intPtr(0),
				PopularityScore:    intPtr(5),
				UrgencyScore:       intPtr(5),
			},
			apiKey:   "trendyol-key",
			expected: http.StatusOK, // 8*30 + 7*25 + 5*20 + 5*15 = 240 + 175 + 100 + 75 = 590
			contains: "590 sent to que",
		},
		{
			name: "New product, Yerel store",
			request: &models.Request{
				ApiKey:             strPtr("yerel-key"),
				ProductName:        strPtr("Random Charger"),
				ProductDescription: strPtr("Generic USB-C Charger"),
				ProductImage:       strPtr("https://example.com/charger.jpg"),
				StoreName:          strPtr("Yerel"),
				Price:              intPtr(100),
				Stock:              intPtr(10),
				PopularityScore:    intPtr(1),
				UrgencyScore:       intPtr(2),
			},
			apiKey:   "yerel-key",
			expected: http.StatusOK, // 3*30 + 7*25 + 1*20 + 2*15 = 90 + 175 + 20 + 30 = 315
			contains: "315 sent to que",
		},
		{
			name: "Big price drop, new merchant",
			request: &models.Request{
				ApiKey:             strPtr("hepsiburada-key"),
				ProductName:        strPtr("IPhone 16"),
				ProductDescription: strPtr("Latest Apple flagship phone"),
				ProductImage:       strPtr("https://example.com/iphone16.jpg"),
				StoreName:          strPtr("Hepsiburada"),
				Price:              intPtr(1050),
				Stock:              intPtr(50),
				PopularityScore:    intPtr(5),
				UrgencyScore:       intPtr(5),
			},
			apiKey:   "hepsiburada-key",
			expected: http.StatusOK, // 7*30 + 7*25 + 5*20 + 5*15 = 210 + 175 + 100 + 75 = 560
			contains: "560 sent to que",
		},		
		{
			name: "Image change only, existing product",
			request: &models.Request{
				ApiKey:             strPtr("amazon-key"),
				ProductName:        strPtr("IPhone 16"),
				ProductDescription: strPtr("Latest Apple flagship phone"),
				ProductImage:       strPtr("https://example.com/iphone16_v2.jpg"), // triggers image change
				StoreName:          strPtr("Amazon"),
				Price:              intPtr(1500),
				Stock:              intPtr(50),
				PopularityScore:    intPtr(5),
				UrgencyScore:       intPtr(5),
			},
			apiKey:   "amazon-key",
			expected: http.StatusOK, // 10*30 + 3*25 + 5*20 + 5*15 = 300 + 75 + 100 + 75 = 550
			contains: "550 sent to que",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			var err error

			if tc.rawBody != nil {
				body = tc.rawBody
			} else {
				body, err = json.Marshal(tc.request)
				if err != nil {
					t.Fatalf("Failed to marshal JSON: %v", err)
				}
			}

			status, response := m.sendRequest(app, t, body, tc.apiKey)

			if status != tc.expected {
				t.Errorf("[%s] Expected status %d, got %d", tc.name, tc.expected, status)
			}
			if !bytes.Contains([]byte(response), []byte(tc.contains)) {
				t.Errorf("[%s] Expected response to contain %q, got %q", tc.name, tc.contains, response)
			}
		})
	}
}
