package registry

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
	"github.com/iomarmochtar/cir-rotator/pkg/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GCRTagsResponse struct {
	Child    []string             `json:"child"`
	Manifest map[string]GCRDigest `json:"manifest"`
	Name     string               `json:"name"`
	Tags     []string             `json:"tags"`
	ErrorsField
}

type GCRDigest struct {
	ImageSizeBytes string   `json:"imageSizeBytes"`
	LayerID        string   `json:"layerId"`
	MediaType      string   `json:"mediaType"`
	Tag            []string `json:"tag"`
	TimeCreatedMs  string   `json:"timeCreatedMs"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
	Name           string
}

type GCR struct {
	host         string
	project      string
	hc           http.IHttpClient
	repositories []Repository
}

func NewGCR(host string, hc http.IHttpClient) (ImageRegistry, error) {
	hostSplt := strings.Split(host, "/")
	if len(hostSplt) == 1 {
		return nil, fmt.Errorf("you must specified parent repository after registry host, eg; asia.gcr.io/parent_repo")
	}
	return &GCR{host: hostSplt[0], project: strings.Join(hostSplt[1:], "/"), hc: hc}, nil
}

// Catalog list of repositorry, recursively follow child attribute
// instead using image registry API spec so we can only use less permission in service account
func (g GCR) Catalog() ([]Repository, error) {
	if err := g.tagList(g.project); err != nil {
		return nil, err
	}
	return g.repositories, nil
}

func (g *GCR) tagList(repository string) (err error) {
	url := fmt.Sprintf("https://%s/v2/%s/tags/list", g.host, repository)
	var jsonBody GCRTagsResponse
	if err = g.hc.GetMarshalReturnObj(url, &jsonBody); err != nil {
		return err
	}

	if len(jsonBody.Errors) > 0 {
		return fmt.Errorf("[%s] [%s]", jsonBody.Errors[0].Code, jsonBody.Errors[0].Message)
	}

	log.Debug().Str("repo", repository).Msg("processing")
	if len(jsonBody.Child) != 0 {
		log.Debug().Msgf("detected %d child(s) in repository %s", len(jsonBody.Child), repository)
		for _, child := range jsonBody.Child {
			nextRepo := fmt.Sprintf("%s/%s", repository, child)
			if err = g.tagList(nextRepo); err != nil {
				return err
			}
		}
	}
	// ignore if it's doesn't has any manifest attached
	if len(jsonBody.Manifest) == 0 {
		log.Debug().Str("repo", repository).Msg("not found any digests found, skipping")
		return nil
	}

	//nolint:prealloc
	var digest []Digest
	for name, gdigest := range jsonBody.Manifest {
		sizeByte, err := strconv.ParseUint(gdigest.ImageSizeBytes, 10, 64)
		if err != nil {
			return errors.Wrap(err, "while converting image size")
		}

		timeCreated, err := h.ConvertTimeStrToUnix(gdigest.TimeCreatedMs)
		if err != nil {
			return errors.Wrap(err, "while parse created time")
		}

		timeUploaded, err := h.ConvertTimeStrToUnix(gdigest.TimeUploadedMs)
		if err != nil {
			return errors.Wrap(err, "while parse uploaded time")
		}

		digest = append(digest, Digest{
			Name:           name,
			ImageSizeBytes: uint(sizeByte),
			Tag:            gdigest.Tag,
			Created:        timeCreated,
			Uploaded:       timeUploaded,
		})
	}
	// name will be combination between host and repo path
	normalizedRepoName := fmt.Sprintf("%s/%s", g.host, repository)
	g.repositories = append(g.repositories, Repository{Name: normalizedRepoName, Digests: digest})
	return nil
}

func (g GCR) Delete(repository Repository) (err error) {
	shortRepoName := strings.TrimPrefix(repository.Name, fmt.Sprintf("%s/", g.host))
	manifestURL := fmt.Sprintf("https://%s/v2/%s/manifests", g.host, shortRepoName)
	for idr := range repository.Digests {
		digest := repository.Digests[idr]
		// delete the related tags
		for idt := range digest.Tag {
			tagURL := fmt.Sprintf("%s/%s", manifestURL, digest.Tag[idt])
			log.Debug().Str("url", tagURL).Msg("deleting tag")
			if err = deleteImage(g.hc, tagURL); err != nil {
				return err
			}
		}

		digestURL := fmt.Sprintf("%s/%s", manifestURL, digest.Name)
		log.Debug().Str("url", digestURL).Msg("deleting digest")
		if err = deleteImage(g.hc, digestURL); err != nil {
			return err
		}
	}

	// delete by digests
	return nil
}

func gcrOauthSource(saData []byte) (oauth2.TokenSource, error) {
	conf, err := google.JWTConfigFromJSON(saData, "https://www.googleapis.com/auth/devstorage.read_write")
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(context.Background()), nil
}
