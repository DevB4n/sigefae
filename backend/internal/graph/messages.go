package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) ListMessages(lastSync time.Time) ([]Message, error) {
	token, err := c.auth.GetToken()
	if err != nil {
		return nil, err
	}

	values := url.Values{}

	values.Set("$orderby", "receivedDateTime asc")
	values.Set("$top", "100") // Adjust the number as needed

	if !lastSync.IsZero() {
		values.Set(
			"$filter",
			fmt.Sprintf(
				"receivedDateTime ge %s",
				lastSync.UTC().Format(time.RFC3339),
			),
		)
	}

	graphURL := fmt.Sprintf(
		"%s/users/%s/messages?%s",
		baseURL,
		c.userEmail,
		values.Encode(),
	)

	req, err := http.NewRequest(http.MethodGet, graphURL, nil)
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
		return nil, fmt.Errorf(
			"graph request failed: %s\n%s",
			resp.Status,
			body,
		)
	}

	var result MessagesResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}