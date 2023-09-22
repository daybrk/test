package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
)

type router struct {
	log *slog.Logger
}

func NewRouter(log *slog.Logger) *router {
	return &router{log: log}
}

func (gr router) EnrichmentAge(name string) (Age, error) {
	response, err := gr.request(name, "https://api.agify.io/")
	if err != nil {
		return Age{}, err
	}

	var ageResp Age
	err = json.NewDecoder(response.Body).Decode(&ageResp)
	if err != nil {
		gr.log.Info("Ошибка", slog.String("err", err.Error()))
		return Age{}, err
	}
	response.Body.Close()

	return ageResp, nil
}

func (gr router) EnrichmentGender(name string) (Gender, error) {
	response, err := gr.request(name, "https://api.genderize.io/")
	if err != nil {
		return Gender{}, err
	}

	var gender Gender
	err = json.NewDecoder(response.Body).Decode(&gender)
	if err != nil {
		gr.log.Info("Ошибка", slog.String("err", err.Error()))
		return Gender{}, err
	}
	response.Body.Close()

	return gender, nil
}

func (gr router) EnrichmentNationality(name string) (Nationality, error) {
	response, err := gr.request(name, "https://api.nationalize.io/")
	if err != nil {
		return Nationality{}, err
	}

	var nationality Nationality
	err = json.NewDecoder(response.Body).Decode(&nationality)
	if err != nil {
		gr.log.Info("Ошибка", slog.String("err", err.Error()))
		return Nationality{}, err
	}
	response.Body.Close()

	return nationality, nil
}

func (gr router) request(value, urlString string) (*http.Response, error) {
	requestURL, err := url.Parse(urlString)
	if err != nil {
		gr.log.Info("Ошибка при парсинге базового URL:", slog.String("err", err.Error()))
		return nil, nil
	}

	query := requestURL.Query()
	query.Set("name", value)
	requestURL.RawQuery = query.Encode()

	response, err := http.Get(requestURL.String())
	if err != nil {
		gr.log.Info("Ошибка при выполнении GET-запроса:", slog.String("err", err.Error()))
		return nil, nil
	}

	return response, err
}
