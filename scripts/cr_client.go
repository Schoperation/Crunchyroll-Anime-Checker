package scripts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type CrunchyRollClient struct {
	accessToken string
	expiresIn   int
}

func NewCrunchyRollClient() (CrunchyRollClient, error) {
	formValues := url.Values{}
	formValues.Set("username", os.Getenv("CR_USERNAME"))
	formValues.Set("password", os.Getenv("CR_PASSWORD"))
	formValues.Set("grant_type", "password")
	formValues.Set("scope", "offline_access")

	request, err := http.NewRequest(http.MethodPost, "https://beta-api.crunchyroll.com/auth/v1/token", bytes.NewBufferString(formValues.Encode()))
	if err != nil {
		return CrunchyRollClient{}, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", os.Getenv("CR_BASIC_KEY")))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Crunchyroll/3.41.1 Android/1.0 okhttp/4.11.0")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return CrunchyRollClient{}, err
	}

	if response.StatusCode >= 400 {
		return CrunchyRollClient{}, fmt.Errorf("could not get new token, error code %d", response.StatusCode)
	}

	var responseStruct struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		Country      string `json:"country"`
		AccountId    string `json:"account_id"`
		ProfileId    string `json:"profile_id"`
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return CrunchyRollClient{}, err
	}

	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return CrunchyRollClient{}, err
	}

	return CrunchyRollClient{
		accessToken: responseStruct.AccessToken,
		expiresIn:   responseStruct.ExpiresIn,
	}, nil
}

func (client CrunchyRollClient) SendReq(path string, queryParams []string) (*http.Response, error) {
	return nil, nil
}
