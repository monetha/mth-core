package tyk

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

	"gitlab.com/monetha/mth-core/http/rest"
)

const (
	createKeyPath = "/tyk/keys/create"
	keyPath       = "/tyk/keys/%s"

	errorBytesLimit = 4000
)

// RestAPI defines methods to access Tyk REST API.
type RestAPI interface {
	CreateKey(ctx context.Context, session *Session) (string, error)
	RetrieveKey(ctx context.Context, keyID string) (*Session, error)
	UpdateKey(ctx context.Context, keyID string, session *Session) error
	DeleteKey(ctx context.Context, keyID string) error
}

// RestAPIClient makes requests to Tyk REST API.
type RestAPIClient struct {
	reqBuilder *rest.Client
	httpClient *http.Client
	authToken  string
}

var _ RestAPI = &RestAPIClient{}

// NewRestAPIClient creates a new Tyk REST API client.
func NewRestAPIClient(httpClient *http.Client, apiURL string, authToken string) *RestAPIClient {
	return &RestAPIClient{
		reqBuilder: rest.NewClient(apiURL),
		httpClient: httpClient,
		authToken:  authToken,
	}
}

// CreateKey creates a new Tyk key using session data and returns key ID.
func (c *RestAPIClient) CreateKey(ctx context.Context, session *Session) (string, error) {
	req, err := c.reqBuilder.NewEndpoint(ctx).Post(createKeyPath).WithHeader("X-Tyk-Authorization", c.authToken).
		WithBody(session).Request()
	if err != nil {
		return "", errorf("failed to build request: %v", err)
	}

	var respObj struct {
		KeyID *string `json:"key"`
	}
	if err = c.do(req, &respObj); err != nil {
		return "", err
	}

	if respObj.KeyID == nil {
		return "", errorf("key ID wasn't included in the response")
	}
	return *respObj.KeyID, nil
}

// RetrieveKey retrieves a key's session data.
func (c *RestAPIClient) RetrieveKey(ctx context.Context, keyID string) (*Session, error) {
	reqPath := fmt.Sprintf(keyPath, keyID)
	req, err := c.reqBuilder.NewEndpoint(ctx).Get(reqPath).WithHeader("X-Tyk-Authorization", c.authToken).Request()
	if err != nil {
		return nil, errorf("failed to build request: %v", err)
	}

	var session Session
	if err = c.do(req, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// UpdateKey updates a key with PUT method, using given session data.
func (c *RestAPIClient) UpdateKey(ctx context.Context, keyID string, session *Session) error {
	reqPath := fmt.Sprintf(keyPath, keyID)
	req, err := c.reqBuilder.NewEndpoint(ctx).Put(reqPath).WithHeader("X-Tyk-Authorization", c.authToken).
		WithBody(session).Request()
	if err != nil {
		return errorf("failed to build request: %v", err)
	}

	return c.do(req, nil)
}

// DeleteKey deletes a key using ID.
func (c *RestAPIClient) DeleteKey(ctx context.Context, keyID string) error {
	reqPath := fmt.Sprintf(keyPath, keyID)
	req, err := c.reqBuilder.NewEndpoint(ctx).Delete(reqPath).WithHeader("X-Tyk-Authorization", c.authToken).Request()
	if err != nil {
		return errorf("failed to build request: %v", err)
	}

	return c.do(req, nil)
}

func (c *RestAPIClient) do(req *http.Request, respObjPtr interface{}) error {
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

// RestAPIStub is a stub which implements RestAPI.
type RestAPIStub struct{}

// NewRestAPIStub creates a new stub.
func NewRestAPIStub() *RestAPIStub {
	return &RestAPIStub{}
}

// CreateKey creates a new Tyk key using session data and returns key ID.
func (c *RestAPIStub) CreateKey(ctx context.Context, session *Session) (string, error) {
	b := make([]byte, 16)
	rand.Read(b)
	h := md5.Sum(b)
	return string(h[:]), nil
}

// RetrieveKey retrieves a key's session data.
func (c *RestAPIStub) RetrieveKey(ctx context.Context, keyID string) (*Session, error) {
	return NewSession(), nil
}

// UpdateKey updates a key with PUT method, using given session data.
func (c *RestAPIStub) UpdateKey(ctx context.Context, keyID string, session *Session) error {
	return nil
}

// DeleteKey deletes a key using ID.
func (c *RestAPIStub) DeleteKey(ctx context.Context, keyID string) error {
	return nil
}
