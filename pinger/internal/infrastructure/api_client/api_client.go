package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Mixturka/DockerLens/pinger/internal/domain/entities"
)

type ApiClient struct {
	BaseUrl string
}

func NewApiClient(baseUrl string) ApiClient {
	return ApiClient{
		BaseUrl: baseUrl,
	}
}

func (ac ApiClient) SavePingResults(rs []entities.PingResult) error {
	baseURL, err := url.Parse(ac.BaseUrl)
	if err != nil {
		return err
	}
	endpoint := baseURL.ResolveReference(&url.URL{Path: "/api/v1/pings"})
	url := endpoint.String()

	data, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save pings on url: %s, status: %d", url, resp.StatusCode)
	}
	return nil
}
