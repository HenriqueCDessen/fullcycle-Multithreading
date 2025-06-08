package models

type Address struct {
	API          string `json:"api"`
	CEP          string `json:"cep"`
	Street       string `json:"street"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	IBGE         string `json:"ibge,omitempty"`
	DDD          string `json:"ddd,omitempty"`
}

type BrasilAPIResponse struct {
	CEP          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type ViaCEPResponse struct {
	CEP          string `json:"cep"`
	Street       string `json:"logradouro"`
	Complement   string `json:"complemento"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	IBGE         string `json:"ibge"`
	DDD          string `json:"ddd"`
}
