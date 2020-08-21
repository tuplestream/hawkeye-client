package hawkeye

import (
	"os"
	"testing"
)

func TestGetEndpoint(t *testing.T) {
	targets := []string{"", "http://incoming.tuplestream.com", "http://somethingelse:8080",
		"https://somethingelse", "https://somethingdifferententirely/subpath"}
	expected := []string{"443", "80", "8080", "443", "443"}

	for i, s := range targets {
		os.Setenv("TUPLESTREAM_HAWKEYE_TARGET", s)
		result := GetHawkeyeTarget()
		if result.Port() != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], result.Port())
		}
	}
	os.Setenv("TUPLESTREAM_HAWKEYE_TARGET", "") // reset
}

func TestConnect(t *testing.T) {
	os.Setenv("TUPLESTREAM_HAWKEYE_TARGET", "https://incoming.tuplestream.com")
	key := os.Getenv("TUPLESTREAM_KEY")
	conn, writer := InitiateConnection("TESTING", key)
	defer conn.Close()
	writer.WriteString("HELLO\n")
	writer.Flush()
}
