package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/henriquedessen/fullcicle-multithreading/clients"
	"github.com/henriquedessen/fullcicle-multithreading/models"
)

func main() {
	cepFlag := flag.String("cep", "", "CEP a ser consultado (com ou sem formatação)")
	flag.Parse()

	var cepArg string
	if len(flag.Args()) > 0 {
		cepArg = flag.Arg(0)
	}

	cepInput := ""
	if *cepFlag != "" {
		cepInput = *cepFlag
	} else if cepArg != "" {
		cepInput = cepArg
	} else {
		fmt.Println("Uso: cepfinder [--cep=CEP] [CEP]")
		fmt.Println("Exemplos:")
		fmt.Println("  cepfinder --cep=14030-430")
		fmt.Println("  cepfinder 14030430")
		fmt.Println("  cepfinder 14030=430")
		os.Exit(1)
	}

	cleanCEP := normalizeCEP(cepInput)

	if len(cleanCEP) != 8 {
		fmt.Println("CEP inválido. Deve conter 8 dígitos.")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	brasilAPI := clients.NewBrasilAPIClient()
	viaCEP := clients.NewViaCEPClient()

	resultChan := make(chan *models.Address, 2)
	errChan := make(chan error, 2)

	go func() {
		address, err := brasilAPI.GetAddress(cleanCEP)
		if err != nil {
			errChan <- fmt.Errorf("BrasilAPI: %v", err)
			return
		}
		resultChan <- address
	}()

	go func() {
		address, err := viaCEP.GetAddress(cleanCEP)
		if err != nil {
			errChan <- fmt.Errorf("ViaCEP: %v", err)
			return
		}
		resultChan <- address
	}()

	select {
	case address := <-resultChan:
		printAddress(address)
	case err := <-errChan:
		select {
		case address := <-resultChan:
			printAddress(address)
		case <-time.After(100 * time.Millisecond):
			log.Fatalf("Erro ao buscar CEP: %v", err)
		}
	case <-ctx.Done():
		log.Fatal("Timeout: nenhuma API respondeu dentro do tempo limite")
	}
}

func normalizeCEP(cep string) string {
	var builder strings.Builder
	for _, c := range cep {
		if c >= '0' && c <= '9' {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

func printAddress(address *models.Address) {
	jsonData, err := json.MarshalIndent(address, "", "  ")
	if err != nil {
		log.Fatalf("Erro ao formatar resposta: %v", err)
	}

	fmt.Printf("Resposta da API: %s\n%s\n", address.API, string(jsonData))
}
