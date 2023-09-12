package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	go request(ctx, "https://cdn.apicep.com/file/apicep", "06550000")
	go request(ctx, "http://viacep.com.br/ws", "06233030")

	time.Sleep(2 * time.Second)
}

func request(ctx context.Context, url string, cep string) {

	if strings.Contains(url, "apicep") {
		cep = maskCEP(cep)
		url = fmt.Sprintf("%s/%s.json", url, cep)
		apiCEP, err := BuscaAPICep(ctx, url)
		if err != nil {
			println(err.Error())
		}
		data, err := json.Marshal(apiCEP)
		if err != nil {
			println(err.Error())
		}
		fmt.Printf("API: %s\n", url)
		fmt.Println("Resultado:", string(data))
		os.Exit(0)

	} else if strings.Contains(url, "viacep") {
		url = fmt.Sprintf("%s/%s/json", url, cep)
		viaCep, err := BuscaViaCep(ctx, url)
		if err != nil {
			println(err.Error())
		}
		data, err := json.Marshal(viaCep)
		if err != nil {
			println(err.Error())
		}
		fmt.Printf("API: %s/\n", url)
		fmt.Println("Resultado:", string(data))
		os.Exit(0)
	}
}

func maskCEP(cep string) string {
	first := cep[0:5]
	second := cep[5:]

	cep = fmt.Sprintf("%s-%s", first, second)

	return cep
}

func BuscaViaCep(ctx context.Context, url string) (*ViaCEP, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func BuscaAPICep(ctx context.Context, url string) (*ApiCEP, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c ApiCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
