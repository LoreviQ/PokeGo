package pokeAPI

import (
	"fmt"
	"testing"
	"time"
)

func TestGetEndpoint(t *testing.T) {
	client := NewClient(5*time.Second, time.Minute)
	cases := []struct {
		next     string
		previous string
		arg      []string
		expected string
	}{
		{
			next:     "",
			previous: "",
			arg:      nil,
			expected: "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
		},
		{
			next:     "https://pokeapi.co/api/v2/location/?offset=40&limit=20",
			previous: "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
			arg:      nil,
			expected: "https://pokeapi.co/api/v2/location/?offset=40&limit=20",
		},
		{
			next:     "https://pokeapi.co/api/v2/location/?offset=40&limit=20",
			previous: "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
			arg:      []string{"-b"},
			expected: "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
		},
		{
			next:     "https://pokeapi.co/api/v2/location/?offset=40&limit=20",
			previous: "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
			arg:      []string{"-back"},
			expected: "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test Case %v", i), func(t *testing.T) {
			actual, err := client.getEndpoint(c.next, c.previous, c.arg)
			if err != nil {
				t.Errorf("Incorrectly raised error %v", err)
			}
			if actual != c.expected {
				t.Errorf("Incorrect endpoint recieved")
			}
		})
	}
	t.Run("Error Test Case", func(t *testing.T) {
		_, err := client.getEndpoint("", "", []string{"-b"})
		if err == nil {
			t.Errorf("Expected error to be raised")
		}
	})
}
