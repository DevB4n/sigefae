package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListMessages() ([]Message, error) {
	token, err := c.auth.GetToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/users/%s/messages",
		baseURL,
		c.userEmail,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("graph request failed: %s\n%s", resp.Status, body)
	}

	var result MessagesResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}