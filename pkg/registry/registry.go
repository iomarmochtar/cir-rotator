package registry

import (
	"errors"
	"fmt"
	"time"

	"regexp"

	"github.com/iomarmochtar/cir-rotator/pkg/http"
	"golang.org/x/oauth2"
)

const (
	GoogleContainerRegistry = "gcr"
)

type (
	oauthTokenSource func(saFilePath string) (oauth2.TokenSource, error)
	registryGen      func(host string, httpClient http.IHttpClient) (ImageRegistry, error)
)

var (
	reGcrMatcher   = regexp.MustCompile(`([a-z]+\.)?(gcr\.io|pkg\.dev)`)
	RegistryMapper = map[string]registryGen{
		GoogleContainerRegistry: NewGCR,
	}
	TokenSourceMapper = map[string]oauthTokenSource{
		GoogleContainerRegistry: gcrOauthSource,
	}
	SupportedContainerRegistryList = []string{
		GoogleContainerRegistry,
	}
)

type ErrorField struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorsField struct {
	Errors []ErrorField `json:"errors"`
}

type Digest struct {
	ImageSizeBytes uint      `json:"size"`
	Tag            []string  `json:"tags"`
	Created        time.Time `json:"created"`
	Uploaded       time.Time `json:""`
	Name           string    `json:"digest"`
}

type Repository struct {
	Name    string   `json:"repository"`
	Digests []Digest `json:"digests"`
}

//go:generate mockgen -destination mock_registry/mock_registry.go -source registry.go ImageRegistry
type ImageRegistry interface {
	Catalog() ([]Repository, error)
	Delete(repo Repository) error
}

// deleteImage shorthand for deleting image
func deleteImage(hc http.IHttpClient, url string) (err error) {
	var errResp ErrorsField
	if err = hc.DeleteMarshalReturnObj(url, &errResp); err != nil {
		return err
	}

	if len(errResp.Errors) != 0 {
		return errors.New(errResp.Errors[0].Message)
	}

	return nil
}

func GetImageRegistryTypeByHostname(host string) (string, error) {
	if reGcrMatcher.Match([]byte(host)) {
		return GoogleContainerRegistry, nil
	}

	return "", fmt.Errorf("unknown matcher registry handler by host %s", host)
}
