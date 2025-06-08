package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/henriquedessen/fullcicle-multithreading/models"
)

type BrasilAPIClient struct {
	BaseURL string
	Timeout time.Duration
}

func NewBrasilAPIClient() *BrasilAPIClient {
	return &BrasilAPIClient{
		BaseURL: "https://brasilapi.com.br/api/cep/v1",
		Timeout: 1 * time.Second,
	}
}

func (c *BrasilAPIClient) GetAddress(cep string) (*models.Address, error) {
	cleanCEP := strings.ReplaceAll(cep, "-", "")

	client := http.Client{Timeout: c.Timeout}
	url := fmt.Sprintf("%s/%s", c.BaseURL, cleanCEP)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("BrasilAPI request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BrasilAPI returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read BrasilAPI response: %v", err)
	}

	var apiResponse models.BrasilAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode BrasilAPI response: %v", err)
	}

	address := &models.Address{
		API:          "BrasilAPI",
		CEP:          apiResponse.CEP,
		Street:       apiResponse.Street,
		Neighborhood: apiResponse.Neighborhood,
		City:         apiResponse.City,
		State:        apiResponse.State,
	}

	return address, nil
}
