package hawkeye

import "testing"

func TestGetEndpoint(t *testing.T) {
	result := GetHawkeyeTarget()
	expected := "incoming.tuplestream.com"
	if result.Host != expected {
		t.Errorf("Expected %s, got %s", expected, result.Host)
	}
}
