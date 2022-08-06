package yolink

// YoLink API Client
// See http://doc.yosmart.com/docs/protocol/openAPIV2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
)

type apiClient struct {
	// User Access Credentials from YoLink app:
	UAID      string
	SecretKey string

	// Set by getToken:
	AccessToken  string
	RefreshToken string

	// Set by getHomeId
	HomeId string

	// internal
	http http.Client
}

type JsonResponse map[string]any
type JsonRequest map[string]any

type Device struct {
	DeviceID   string
	DeviceUUID string
	Token      string
	Name       string
	Type       string
}

var ErrInvalidFormat = errors.New("invalid format")

func NewAPIClient(uaid, secretKey string) (*apiClient, error) {
	client := apiClient{UAID: uaid, SecretKey: secretKey}
	client.http = http.Client{Timeout: time.Duration(5) * time.Second}

	_, err := client.getToken(false)
	if err != nil {
		return nil, err
	}

	// Home ID is required for MQTT connection:
	_, err = client.GetHomeId()
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *apiClient) GetHomeId() (string, error) {
	if c.HomeId != "" {
		// use previous result
		return c.HomeId, nil
	}

	response, err := c.apiRequest(
		JsonRequest{
			"method": "Home.getGeneralInfo",
		},
	)
	if err != nil {
		return "", err
	}

	var ok bool
	data, ok := response["data"].(map[string]any)
	if !ok {
		return "", ErrInvalidFormat
	}
	c.HomeId, ok = data["id"].(string)
	if !ok {
		return "", ErrInvalidFormat
	}
	return c.HomeId, nil
}

func (c *apiClient) GetDevices() ([]Device, error) {
	devices := make([]Device, 0)

	response, err := c.apiRequest(
		JsonRequest{
			"method": "Home.getDeviceList",
		},
	)
	if err != nil {
		return devices, err
	}

	var ok bool
	data, ok := response["data"].(map[string]any)
	if !ok {
		return devices, ErrInvalidFormat
	}
	deviceList, ok := data["devices"].([]any)
	if !ok {
		return devices, ErrInvalidFormat
	}
	for _, md := range deviceList {
		dev := Device{}
		mapstructure.Decode(md, &dev)
		devices = append(devices, dev)
	}

	return devices, nil
}

// Private methods:

// Requests access token from UAID and Secret Key
func (c *apiClient) getToken(refresh bool) (JsonResponse, error) {
	r := make(JsonResponse)
	form := url.Values{}
	if refresh {
		form.Set("grant_type", "refresh_token")
		form.Set("client_id", c.UAID)
		form.Set("refresh_token", c.RefreshToken)
	} else {
		form.Set("grant_type", "client_credentials")
		form.Set("client_id", c.UAID)
		form.Set("client_secret", c.SecretKey)
	}
	resp, err := c.http.PostForm(serverURL(TOKEN_PATH), form)
	if err != nil {
		return r, err
	}

	if resp.StatusCode != 200 {
		return r, fmt.Errorf("http status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return r, err
	}

	if state, ok := r["state"]; ok {
		if state == ERROR_STATE {
			if msg, ok := r["msg"].(string); ok {
				return r, errors.New(msg)
			}
			return r, errors.New("unknown token error")
		}
	}

	var ok bool
	c.AccessToken, ok = r["access_token"].(string)
	if !ok {
		return r, ErrInvalidFormat
	}
	c.RefreshToken, ok = r["refresh_token"].(string)
	if !ok {
		return r, ErrInvalidFormat
	}

	return r, nil
}

// Makes a generic API JSON request and returns response
func (c *apiClient) apiRequest(body JsonRequest) (JsonResponse, error) {
	r := make(JsonResponse)
	js, err := json.Marshal(body)
	if err != nil {
		return r, err
	}
	req, err := http.NewRequest("POST", serverURL(API_PATH), bytes.NewBuffer(js))
	if err != nil {
		return r, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.http.Do(req)
	if err != nil {
		return r, err
	}

	if resp.StatusCode != 200 {
		return r, fmt.Errorf("http status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return r, err
	}

	if code, ok := r["code"].(string); ok {
		// API is not RESTful, does not use HTTP status codes for all errors
		if code != SUCCESS_CODE {
			if desc, ok := r["desc"].(string); ok {
				return r, fmt.Errorf("error #%s: %s", code, desc)
			} else {
				return r, fmt.Errorf("error #%s: unknown", code)
			}
		}
	}

	return r, nil
}

func serverURL(path string) string {
	return fmt.Sprintf("https://%s%s", API_HOST, path)
}
