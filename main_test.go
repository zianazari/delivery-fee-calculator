package main

import "testing"

func TestCalculateDeliveryFee(t *testing.T) {
	tests := []struct {
		name     string
		request  CalculateFeeRequest
		expected int
	}{
		{
			name: "Free delivery for cart value >= 200€",
			request: CalculateFeeRequest{
				CartValue:        20000,
				DeliveryDistance: 1500,
				NumberOfItems:    3,
				Time:             "2026-02-05T14:00:00Z",
			},
			expected: 0, // Free delivery for cart value >= 200€
		},
		{
			name: "Basic delivery fee calculation",
			request: CalculateFeeRequest{
				CartValue:        1000,
				DeliveryDistance: 1500,
				NumberOfItems:    3,
				Time:             "2026-02-05T14:00:00Z",
			},
			expected: 300, // 2€ for first 1000m + 1€ for additional 500m
		},
		{
			name: "Basic delivery fee calculation 2",
			request: CalculateFeeRequest{
				CartValue:        1000,
				DeliveryDistance: 1501,
				NumberOfItems:    3,
				Time:             "2026-02-05T14:00:00Z",
			},
			expected: 400, // 2€ for first 1000m + 1€ for additional first 500m + 1€ for additional second 500m
		},
		{
			name: "Small order surcharge",
			request: CalculateFeeRequest{
				CartValue:        800,
				DeliveryDistance: 800,
				NumberOfItems:    3,
				Time:             "2026-02-05T14:00:00Z",
			},
			expected: 400, // 2€ small order surcharge + 2€ delivery fee
		},
		{
			name: "Surcharge for number of items",
			request: CalculateFeeRequest{
				CartValue:        1500,
				DeliveryDistance: 800,
				NumberOfItems:    6,
				Time:             "2026-02-05T14:00:00Z",
			},
			expected: 300, // 2€ delivery + 1€ (50¢ per item above 4)
		},
		{
			name: "Friday rush hour multiplier",
			request: CalculateFeeRequest{
				CartValue:        1500,
				DeliveryDistance: 1500,
				NumberOfItems:    3,
				Time:             "2026-02-06T16:00:00Z",
			},
			expected: 360, // (2€ + 1€) * 1.2 = 3.6€
		},
		{
			name: "Maximum delivery fee cap",
			request: CalculateFeeRequest{
				CartValue:        500,
				DeliveryDistance: 10000,
				NumberOfItems:    20,
				Time:             "2026-02-05T14:00:00Z",
			},
			expected: 1500, // capped at 15€
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee := CalculateDeliveryFee(&tt.request)
			if fee != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, fee)
			}
		})
	}
}
