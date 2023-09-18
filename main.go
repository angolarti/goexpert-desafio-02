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

type CepType interface {
	*ViaCEP | *ApiCEP
}

type CEP[T CepType] struct {
	Body T
}

func RefCEPType[T CepType](cep T) *CEP[T] {
	return &CEP[T]{
		Body: cep,
	}
}

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

	select {
	case <-time.After(time.Second * 1):
		go request(ctx, "https://cdn.apicep.com/file/apicep", "06550000", &ApiCEP{})
		go request(ctx, "http://viacep.com.br/ws", "06233030", &ViaCEP{})
		time.Sleep(10 * time.Second)
		return
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println(err.Error())
	}

}

func request[T CepType](ctx context.Context, url string, cep string, api T) {

	if strings.Contains(url, "apicep") {
		cep = maskCEP(cep)
		url = fmt.Sprintf("%s/%s.json", url, cep)
	} else if strings.Contains(url, "viacep") {
		url = fmt.Sprintf("%s/%s/json", url, cep)
	}
	response, err := BauscaCEP(ctx, "GET", url, api)

	if err != nil {
		println(err.Error())
	}

	data, err := json.Marshal(response)
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("API: %s\n", url)
	fmt.Println("Resultado:", string(data))
	os.Exit(0)

}

func maskCEP(cep string) string {
	first := cep[0:5]
	second := cep[5:]

	cep = fmt.Sprintf("%s-%s", first, second)

	return cep
}

func BauscaCEP[T CepType](ctx context.Context, method string, url string, object T) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)

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

	c := RefCEPType(object)

	err = json.Unmarshal(body, c.Body)
	if err != nil {
		return nil, err
	}

	return &c.Body, nil
}
