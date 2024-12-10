package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StaticDhcpHost struct {
	MacAddress string
	IPAddress  string
	HostName   string
}

type errorJSON struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type Client interface {
	CreateStaticDhcpHost(host StaticDhcpHost) (*StaticDhcpHost, error)
	ReadStaticDhcpHost(macAddress string) (*StaticDhcpHost, error)
	UpdateStaticDhcpHost(host StaticDhcpHost) (*StaticDhcpHost, error)
	DeleteStaticDhcpHost(macAddress string) (*StaticDhcpHost, error)
}

func New(apiUrl string, token string) Client {
	return &dnsmasqManagerClient{
		httpClient: http.DefaultClient,
		apiUrl:     apiUrl,
		jwtToken:   token,
	}
}

type dnsmasqManagerClient struct {
	httpClient *http.Client
	apiUrl     string
	jwtToken   string
}

func (c *dnsmasqManagerClient) CreateStaticDhcpHost(host StaticDhcpHost) (*StaticDhcpHost, error) {
	return c.staticDhcpHostRequestWithBody(http.MethodPost, host)
}

func (c *dnsmasqManagerClient) ReadStaticDhcpHost(macAddress string) (*StaticDhcpHost, error) {
	return c.staticDhcpHostRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/static/host?mac=%s", c.apiUrl, macAddress),
		nil,
		http.StatusOK)
}

func (c *dnsmasqManagerClient) UpdateStaticDhcpHost(host StaticDhcpHost) (*StaticDhcpHost, error) {
	return c.staticDhcpHostRequestWithBody(http.MethodPut, host)
}

func (c *dnsmasqManagerClient) DeleteStaticDhcpHost(macAddress string) (*StaticDhcpHost, error) {
	return c.staticDhcpHostRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/static/host?mac=%s", c.apiUrl, macAddress),
		nil,
		http.StatusOK)
}

func (c *dnsmasqManagerClient) staticDhcpHostRequestWithBody(httpMethod string, host StaticDhcpHost) (*StaticDhcpHost, error) {
	body, err := json.Marshal(&host)
	if err != nil {
		return nil, err
	}

	return c.staticDhcpHostRequest(
		httpMethod,
		fmt.Sprintf("%s/api/v1/static/host", c.apiUrl),
		strings.NewReader(string(body)),
		http.StatusCreated)
}

func (c *dnsmasqManagerClient) staticDhcpHostRequest(httpMethod string, url string, body io.Reader, successStatus int) (*StaticDhcpHost, error) {
	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}

	response_body, err := c.doRequest(request, successStatus)
	if err != nil {
		return nil, err
	}

	host := StaticDhcpHost{}
	err = json.Unmarshal(response_body, &host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (c *dnsmasqManagerClient) doRequest(request *http.Request, successStatus int) ([]byte, error) {
	if c.jwtToken != "" {
		request.Header.Set("Authorization", "Bearer "+c.jwtToken)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != successStatus {
		response_error := errorJSON{}
		err = json.Unmarshal(body, &response_error)
		if err != nil {
			return nil, fmt.Errorf("Status: %d\nBody: %s", response.StatusCode, body)
		} else {
			return nil, fmt.Errorf("%s\n\n%s", response_error.Message, response_error.Details)
		}
	}

	return body, err
}
