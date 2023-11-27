package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"
)

// CrunchyrollClient represents a client capable of sending HTTP requests to Crunchyroll.
type CrunchyrollClient struct {
	credFilePath string
	listsPath    string
	accessToken  string
	expiresIn    int
	locale       Locale
	lastLogin    time.Time
}

// Used for loading creds from the cred file.
type credsFromFile struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	BasicAuthKey string `json:"basic_auth_key"`
}

func NewCrunchyrollClient(credFilePath string, listsPath string, locale Locale) (CrunchyrollClient, error) {
	client := CrunchyrollClient{
		credFilePath: credFilePath,
		listsPath:    listsPath,
		locale:       locale,
	}

	err := client.Login()
	if err != nil {
		return CrunchyrollClient{}, err
	}

	return client, nil
}

func (client *CrunchyrollClient) Login() error {
	credFile, err := openCredFile(client.credFilePath)
	if err != nil {
		return err
	}

	formValues := url.Values{}
	formValues.Set("username", credFile.Username)
	formValues.Set("password", credFile.Password)
	formValues.Set("grant_type", "password")
	formValues.Set("scope", "offline_access")

	request, err := http.NewRequest(http.MethodPost, "https://beta-api.crunchyroll.com/auth/v1/token", bytes.NewBufferString(formValues.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", credFile.BasicAuthKey))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Crunchyroll/3.41.1 Android/1.0 okhttp/4.11.0")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("could not get new token, error code %d", response.StatusCode)
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
		return err
	}

	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return err
	}

	client.accessToken = responseStruct.AccessToken
	client.expiresIn = responseStruct.ExpiresIn
	client.lastLogin = time.Now()
	return nil
}

func (client *CrunchyrollClient) Get(path string, responseStruct any) error {
	return client.GetWithQueryParams(path, responseStruct, map[string]string{
		"locale": client.locale.Name(),
	})
}

func (client *CrunchyrollClient) GetWithQueryParams(path string, responseStruct any, queryParams map[string]string) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://beta-api.crunchyroll.com/%s", path), nil)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.accessToken))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Crunchyroll/3.41.1 Android/1.0 okhttp/4.11.0")

	values := url.Values{}
	for key, value := range queryParams {
		values.Set(key, value)
	}

	request.URL.RawQuery = values.Encode()

	// Re-login if we believe the token is expired about now.
	if time.Since(client.lastLogin).Seconds() > math.Abs(float64(client.expiresIn)-100) {
		err = client.Login()
		if err != nil {
			return err
		}
	}

	backOffSchedule := []time.Duration{
		1 * time.Second,
		5 * time.Second,
		10 * time.Second,
	}

	var response *http.Response
	for _, backoff := range backOffSchedule {
		response, err = http.DefaultClient.Do(request)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusOK {
			break
		} else if response.StatusCode == http.StatusGatewayTimeout || response.StatusCode == http.StatusServiceUnavailable || response.StatusCode == http.StatusRequestTimeout {
			time.Sleep(backoff)
			continue
		} else {
			return fmt.Errorf("failed to get response from crunchyroll; status code %d", response.StatusCode)
		}
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseStruct)
	if err != nil {
		return err
	}

	return nil
}

func openCredFile(credFilePath string) (credsFromFile, error) {
	file, err := os.Open(credFilePath)
	if err != nil {
		return credsFromFile{}, err
	}

	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return credsFromFile{}, err
	}

	var creds credsFromFile
	err = json.Unmarshal(bytes, &creds)
	if err != nil {
		return credsFromFile{}, err
	}

	return creds, nil
}
