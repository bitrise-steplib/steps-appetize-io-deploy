package appetize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
)

const apiEndPoint = "@api.appetize.io/v1/apps"

// Client ...
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
	appPath    string
	artifact   Artifact
	publicKey  string
}

// Response ...
type Response struct {
	PublicKey      string        `json:"publicKey"`
	PrivateKey     string        `json:"privateKey"`
	Updated        time.Time     `json:"updated"`
	Email          string        `json:"email"`
	Platform       string        `json:"platform"`
	VersionCode    int           `json:"versionCode"`
	Created        time.Time     `json:"created"`
	Architectures  []interface{} `json:"architectures"`
	AppPermissions struct {
	} `json:"appPermissions"`
	PublicURL string `json:"publicURL"`
	AppURL    string `json:"appURL"`
	ManageURL string `json:"manageURL"`
}

// -------------------------------------
// -- Public methods

// NewClient ...
func NewClient(token string, appPath string, artifact Artifact, publicKey string) *Client {
	baseURL := baseURL(token, appPath, publicKey)
	return &Client{
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
		appPath:    appPath,
		artifact:   artifact,
	}
}

// DirectFileUpload ...
func (client *Client) DirectFileUpload() (Response, error) {
	request, err := createRequest(client.baseURL, client.appPath, string(client.artifact.Platform()))
	if err != nil {
		return Response{}, err
	}

	var resp Response
	_, err = client.performRequest(request, &resp)
	return resp, err
}

func createRequest(url, appPath, platform string) (*http.Request, error) {
	f, err := os.Open(appPath)
	if err != nil {
		log.Errorf("%s", err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Errorf("%s", err)
	}

	var b bytes.Buffer
	mpw := multipart.NewWriter(&b)

	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Errorf("failed to close file error: %s", cerr)
		}
	}()

	// file
	{
		part, err := mpw.CreateFormFile("file", fi.Name())
		if err != nil {
			log.Errorf("%s", err)
		}

		_, err = io.Copy(part, f)
		if err != nil {
			log.Errorf("%s", err)
		}
	}

	// platform
	{
		field, err := mpw.CreateFormField("platform")
		if err != nil {
			log.Errorf("%s", err)
		}

		_, err = io.Copy(field, strings.NewReader(platform))
		if err != nil {
			log.Errorf("%s", err)
		}
	}

	if err = mpw.Close(); err != nil {
		log.Errorf("%s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mpw.FormDataContentType())
	return req, nil
}

// -------------------------------------
// -- Private methods

func (client *Client) performRequest(req *http.Request, requestResponse interface{}) ([]byte, error) {
	response, err := client.httpClient.Do(req)
	if err != nil {
		// On error, any Response can be ignored
		return nil, fmt.Errorf("failed to perform request, error: %s", err)
	}

	// The client must close the response body when finished with it
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Warnf("Failed to close response body, error: %s", cerr)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body, error: %s", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusMultipleChoices {
		return nil, fmt.Errorf("Response status: %d - Body: %s", response.StatusCode, string(body))
	}

	// Parse JSON body
	if requestResponse != nil {
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}
	}
	return body, nil
}

func baseURL(token, appPath, publicKey string) string {
	baseURL := token + apiEndPoint

	if publicKey != "" {
		baseURL = path.Join(baseURL, publicKey)
	}
	return "https://" + baseURL
}
