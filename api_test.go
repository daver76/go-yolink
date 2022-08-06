package yolink

import (
	"os"
	"testing"
	"time"
)

// Integration test: hits YoLink API with credentials from the environment
func TestNewClient(t *testing.T) {
	uaid := os.Getenv("YOLINK_UAID")
	if uaid == "" {
		t.Fatalf("YOLINK_UAID not set")
	}
	secret_key := os.Getenv("YOLINK_SECRET_KEY")
	if secret_key == "" {
		t.Fatalf("YOLINK_SECRET_KEY not set")
	}

	client, err := NewAPIClient(uaid, secret_key)
	if err != nil {
		t.Fatal(err)
	}
	if client.AccessToken == "" {
		t.Errorf("access_token invalid: %s", client.AccessToken)
	}
	if client.RefreshToken == "" {
		t.Errorf("refresh_token invalid: %s", client.RefreshToken)
	}
	if client.HomeId == "" {
		t.Errorf("home_id invalid: %s", client.HomeId)
	}

	// Retrieve devices
	devices, err := client.GetDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		t.Errorf("no devices found")
	}

	// Refresh access token:
	prevAccessToken := client.AccessToken
	prevRefreshToken := client.RefreshToken
	time.Sleep(2 * time.Second) // needed, otherwise we may get the same token back
	_, err = client.getToken(true)
	if err != nil {
		t.Fatal(err)
	}
	if prevAccessToken == client.AccessToken {
		t.Fatalf("Unable to refresh access token: %s = %s", prevAccessToken, client.AccessToken)
	}
	if prevRefreshToken == client.RefreshToken {
		t.Fatalf("Unable to refresh refresh token")
	}
}
