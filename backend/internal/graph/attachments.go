package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (c *Client) Attachments(messageID string) ([]Attachment, error) {
	token, err := c.auth.GetToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/users/%s/messages/%s/attachments",
		baseURL,
		c.userEmail,
		messageID,
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

	var result AttachmentsResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}

func (c *Client) DownloadAttachments(messageID string) error {
	attachments, err := c.Attachments(messageID)
	if err != nil {
		return err
	}

	if len(attachments) == 0 {
		return nil
	}

	if err := os.MkdirAll("downloads", 0755); err != nil {
		return err
	}

	for _, att := range attachments {
		c.fileCounter++

		filename := fmt.Sprintf(
			"archivo%d%s",
			c.fileCounter,
			filepath.Ext(att.Name),
		)

		path := filepath.Join("downloads", filename)

		if err := os.WriteFile(path, att.ContentBytes, 0644); err != nil {
			return err
		}

		fmt.Printf("✓ Guardado: %s\n", path)
	}

	return nil
}