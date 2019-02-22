package tyk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gitlab.com/monetha/mth-core/http/rest"
)

const (
	tykDashboardAPIBaseURL = "https://admin.cloud.tyk.io"
	createSessionPath      = "/api/keys"
	apiKeyPath             = "/api/apis/%s/keys/%s"

	errorBytesLimit = 4000
)

// Client makes requests to Tyk Dashboard API.
type Client struct {
	reqBuilder *rest.Client
	httpClient *http.Client
	authToken  string
}

// NewClient creates a new Tyk Dashboard API client.
func NewClient(httpClient *http.Client, authToken string) *Client {
	return &Client{
		reqBuilder: rest.NewClient(tykDashboardAPIBaseURL),
		httpClient: httpClient,
		authToken:  authToken,
	}
}

// CreateSession creates a new Tyk session and returns key ID.
func (c *Client) CreateSession(ctx context.Context, session *Session) (string, error) {
	req, err := c.reqBuilder.NewEndpoint(ctx).Post(createSessionPath).WithAuth(&c.authToken).WithBody(session).Request()
	if err != nil {
		return "", errorf("failed to build request: %v", err)
	}

	var respObj struct {
		KeyID *string `json:"key_id"`
	}
	if err = c.do(req, &respObj); err != nil {
		return "", err
	}

	if respObj.KeyID == nil {
		return "", errorf("key ID wasn't included in the response")
	}
	return *respObj.KeyID, nil
}

// RetrieveSession retrieves a session.
func (c *Client) RetrieveSession(ctx context.Context, apiID, keyID string) (*Session, error) {
	reqPath := fmt.Sprintf(apiKeyPath, apiID, keyID)
	req, err := c.reqBuilder.NewEndpoint(ctx).Get(reqPath).WithAuth(&c.authToken).Request()
	if err != nil {
		return nil, errorf("failed to build request: %v", err)
	}

	var respObj struct {
		Data Session `json:"data"`
	}
	if err = c.do(req, &respObj); err != nil {
		return nil, err
	}

	return &respObj.Data, nil
}

// UpdateSession updates a session with PUT method.
func (c *Client) UpdateSession(ctx context.Context, apiID, keyID string, session *Session) error {
	reqPath := fmt.Sprintf(apiKeyPath, apiID, keyID)
	req, err := c.reqBuilder.NewEndpoint(ctx).Put(reqPath).WithAuth(&c.authToken).WithBody(session).Request()
	if err != nil {
		return errorf("failed to build request: %v", err)
	}

	return c.do(req, nil)
}

// DeleteSession deletes a session using a relevant API ID.
func (c *Client) DeleteSession(ctx context.Context, apiID, keyID string) error {
	reqPath := fmt.Sprintf(apiKeyPath, apiID, keyID)
	req, err := c.reqBuilder.NewEndpoint(ctx).Delete(reqPath).WithAuth(&c.authToken).Request()
	if err != nil {
		return errorf("failed to build request: %v", err)
	}

	return c.do(req, nil)
}

func (c *Client) do(req *http.Request, respObjPtr interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errorify(resp.Body)
	}

	if respObjPtr != nil {
		if err = json.NewDecoder(resp.Body).Decode(respObjPtr); err != nil {
			return errorf("failed to decode successful response body: %v", err)
		}
	}

	return nil
}

func errorify(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errorf("failed to read bytes to errorify: %v", err)
	}
	if len(b) > errorBytesLimit {
		b = b[:errorBytesLimit]
	}
	return errors.New(string(b))
}

func errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("tyk: "+msg, args...)
}
