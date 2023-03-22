package tyk

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"path"

	"github.com/monetha/mth-core/http/rest"
)

const (
	keyPath = "/tyk/keys/"

	errorBytesLimit = 4000
)

// ErrBadRequest is returned when request was incorrect
var ErrBadRequest = errors.New("bad request")

// ErrUnauthorized is returned when authentication information is missing
var ErrUnauthorized = errors.New("authorization required")

// ErrForbidden is returned when access is denied
var ErrForbidden = errors.New("forbidden")

// ErrNotFound is returned when a key is not found.
var ErrNotFound = errors.New("key not found")

// RestAPI defines methods to access Tyk REST API.
type RestAPI interface {
	CreateKey(ctx context.Context, session *Session, key string) (string, error)
	RetrieveKey(ctx context.Context, key string) (*Session, error)
	UpdateKey(ctx context.Context, key string, session *Session) error
	DeleteKey(ctx context.Context, key string) error
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
func (c *RestAPIClient) CreateKey(ctx context.Context, session *Session, key string) (string, error) {
	reqPath := path.Join(keyPath, key)
	req, err := c.reqBuilder.NewEndpoint(ctx).Post(reqPath).WithHeader("X-Tyk-Authorization", c.authToken).
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
func (c *RestAPIClient) RetrieveKey(ctx context.Context, key string) (*Session, error) {
	reqPath := path.Join(keyPath, key)
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
func (c *RestAPIClient) UpdateKey(ctx context.Context, key string, session *Session) error {
	reqPath := path.Join(keyPath, key)
	req, err := c.reqBuilder.NewEndpoint(ctx).Put(reqPath).WithHeader("X-Tyk-Authorization", c.authToken).
		WithBody(session).Request()
	if err != nil {
		return errorf("failed to build request: %v", err)
	}

	return c.do(req, nil)
}

// DeleteKey deletes a key using ID.
func (c *RestAPIClient) DeleteKey(ctx context.Context, key string) error {
	reqPath := path.Join(keyPath, key)
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

	code := resp.StatusCode

	if code >= 400 {
		switch code {
		case 400:
			return ErrBadRequest
		case 401:
			return ErrUnauthorized
		case 403:
			return ErrForbidden
		case 404:
			return ErrNotFound
		default:
			return errorify(resp.Body)
		}
	}

	if respObjPtr != nil {
		if err = json.NewDecoder(resp.Body).Decode(respObjPtr); err != nil {
			return errorf("failed to decode successful response body: %v", err)
		}
	}

	return nil
}

func errorify(r io.Reader) error {
	b, err := io.ReadAll(r)
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

var _ RestAPI = &RestAPIStub{}

// NewRestAPIStub creates a new stub.
func NewRestAPIStub() *RestAPIStub {
	return &RestAPIStub{}
}

// CreateKey creates a new Tyk key using session data and returns key ID.
func (c *RestAPIStub) CreateKey(ctx context.Context, session *Session, key string) (string, error) {
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
