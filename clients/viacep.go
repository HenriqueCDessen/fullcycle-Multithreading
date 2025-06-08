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

type ViaCEPClient struct {
	BaseURL string
	Timeout time.Duration
}

func NewViaCEPClient() *ViaCEPClient {
	return &ViaCEPClient{
		BaseURL: "http://viacep.com.br/ws",
		Timeout: 1 * time.Second,
	}
}

func (c *ViaCEPClient) GetAddress(cep string) (*models.Address, error) {
	cleanCEP := strings.ReplaceAll(cep, "-", "")

	client := http.Client{Timeout: c.Timeout}
	url := fmt.Sprintf("%s/%s/json/", c.BaseURL, cleanCEP)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ViaCEP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read ViaCEP response: %v", err)
	}

	if strings.Contains(string(body), `"erro": true`) {
		return nil, fmt.Errorf("CEP não encontrado na ViaCEP")
	}

	var apiResponse models.ViaCEPResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode ViaCEP response: %v", err)
	}

	address := &models.Address{
		API:          "ViaCEP",
		CEP:          apiResponse.CEP,
		Street:       apiResponse.Street,
		Complement:   apiResponse.Complement,
		Neighborhood: apiResponse.Neighborhood,
		City:         apiResponse.City,
		State:        apiResponse.State,
		IBGE:         apiResponse.IBGE,
		DDD:          apiResponse.DDD,
	}

	return address, nil
}
