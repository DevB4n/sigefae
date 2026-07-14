package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Auth struct {
	clientID     string
	clientSecret string
	
	tenantID     string

	token     string
	expiresAt time.Time
}

type tokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func NewAuth(clientID, clientSecret, tenantID string) *Auth {
	return &Auth{
		clientID:     clientID,
		clientSecret: clientSecret,
		tenantID:     tenantID,
	}
}

func (a *Auth) GetToken() (string, error) {
	// Si el token aún es válido, reutilizarlo
	if a.token != "" && time.Now().Before(a.expiresAt) {
		return a.token, nil
	}

	form := url.Values{}
	form.Set("client_id", a.clientID)
	form.Set("client_secret", a.clientSecret)
	form.Set("scope", "https://graph.microsoft.com/.default")
	form.Set("grant_type", "client_credentials")

	tokenURL := fmt.Sprintf(
		"https://login.microsoftonline.com/%s/oauth2/v2.0/token",
		a.tenantID,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		tokenURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("authentication failed: %s\n%s", resp.Status, body)
	}

	var tokenResp tokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	a.token = tokenResp.AccessToken

	// Renovar un minuto antes de expirar
	a.expiresAt = time.Now().Add(
		time.Duration(tokenResp.ExpiresIn-60) * time.Second,
	)

	return a.token, nil
}
