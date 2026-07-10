package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	zipExtractor "sigefae/internal/zip"
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

func (c *Client) DownloadAttachments(message Message) error {
	attachments, err := c.Attachments(message.ID)
	if err != nil {
		return err
	}

	if len(attachments) == 0 {
		return nil
	}

	dir := filepath.Join("downloads", message.ID)

	// Si ya existe la carpeta, este correo ya fue procesado.
	if _, err := os.Stat(dir); err == nil {
		fmt.Printf("✓ Correo ya descargado: %s\n", message.Subject)
		return nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	metadata := MailMetadata{
		ID:               message.ID,
		Subject:          message.Subject,
		ReceivedDateTime: message.ReceivedDateTime,
		From:             message.From.EmailAddress.Address,
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		filepath.Join(dir, "correo.json"),
		data,
		0644,
	); err != nil {
		return err
	}

	extractor := zipExtractor.NewExtractor()

	for _, att := range attachments {
		path := filepath.Join(dir, filepath.Base(att.Name))

		if err := os.WriteFile(path, att.ContentBytes, 0644); err != nil {
			return err
		}

		fmt.Printf("✓ Guardado: %s\n", path)

		// Si el adjunto es un ZIP, extraerlo automáticamente.
			if strings.EqualFold(filepath.Ext(att.Name), ".zip") {
			destination := filepath.Join(dir, "extracted")

			if err := extractor.Extract(path, destination); err != nil {
				return err
			}

			fmt.Printf("✓ ZIP extraído en: %s\n", destination)
		}
	}

	return nil
}