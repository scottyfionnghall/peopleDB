package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GenderResp struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type AgeResp struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type NationalityResp struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Info interface {
	GenderResp | AgeResp | NationalityResp
}

func get(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		err := errors.New("third party api not available")
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func neededInfo[V Info](url string, info V) (V, error) {
	data, err := get(url)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

func populateInfo(name string) (int, string, string, error) {

	var (
		age         AgeResp
		gender      GenderResp
		nationality NationalityResp
	)

	age_url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	gender_url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	nationality_url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)

	age, err := neededInfo(age_url, age)
	if err != nil {
		return -1, "", "", err
	}

	gender, err = neededInfo(gender_url, gender)
	if err != nil {
		return -1, "", "", err
	}

	nationality, err = neededInfo(nationality_url, nationality)
	if err != nil {
		return -1, "", "", err
	}

	return age.Age, gender.Gender, nationality.Country[0].CountryId, nil
}
