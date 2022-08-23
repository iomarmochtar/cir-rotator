package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	fl "github.com/iomarmochtar/cir-rotator/pkg/filter"
	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
	http "github.com/iomarmochtar/cir-rotator/pkg/http"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
)

//go:generate mockgen -destination mock_config/mock_config.go -source config.go IConfig
type IConfig interface {
	Username() string
	Password() string
	SkipList() []string
	IsDryRun() bool
	Host() string
	ImageRegistry() reg.ImageRegistry
	ExcludeEngine() fl.IFilterEngine
	IncludeEngine() fl.IFilterEngine
	HTTPClient() http.IHttpClient
	RepositoryList() []reg.Repository
	Init() error
}

type Config struct {
	RegUsername        string
	RegPassword        string
	ServiceAccountPath string
	RegistryHost       string
	RegistryType       string
	SkipListPath       string
	RepoListPath       string
	DryRun             bool
	ExcludeFilters     []string
	IncludeFilters     []string
	AllowInsecure      bool
	JWExpirySecond     uint

	excludeEngine fl.IFilterEngine
	includeEngine fl.IFilterEngine
	imageReg      reg.ImageRegistry
	httpClient    http.IHttpClient
	skipList      []string
	repositories  []reg.Repository
}

// Init is validating inputs and setups some dependencies for the application to run
func (c *Config) Init() (err error) {
	// registry type
	if err = c.initRegType(); err != nil {
		return err
	}

	// http client
	if err = c.initHTTPClient(); err != nil {
		return err
	}

	// setup image registry
	if err = c.initImageReg(); err != nil {
		return err
	}

	// skip list that will be used in delete actions
	if err = c.initSkipList(); err != nil {
		return err
	}

	// setup filters (include & exclude)
	if err = c.initFilters(); err != nil {
		return err
	}

	// repositories list
	if err = c.initRepositoryList(); err != nil {
		return err
	}

	return nil
}

func (c Config) SkipList() []string {
	return c.skipList
}

func (c Config) RepositoryList() []reg.Repository {
	return c.repositories
}

func (c Config) IsDryRun() bool {
	return c.DryRun
}

func (c Config) Username() string {
	return c.RegUsername
}

func (c Config) Password() string {
	return c.RegPassword
}

func (c Config) Host() string {
	return c.RegistryHost
}

func (c Config) HTTPClient() http.IHttpClient {
	return c.httpClient
}

func (c Config) ExcludeEngine() fl.IFilterEngine {
	return c.excludeEngine
}

func (c Config) IncludeEngine() fl.IFilterEngine {
	return c.includeEngine
}

// ImageRegistry get related image registry based on known host if not mentioned
// directly for image host name
func (c Config) ImageRegistry() reg.ImageRegistry {
	return c.imageReg
}

func (c *Config) initRegType() (err error) {
	if c.Host() == "" {
		return fmt.Errorf("registry host is required")
	}
	// if registry type was not mentioned then try do determinte it by hostname
	if c.RegistryType == "" {
		if c.RegistryType, err = reg.GetImageRegistryByHostname(c.Host()); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) initImageReg() (err error) {
	imageRegFn := reg.RegistryMapper[c.RegistryType]
	if imageRegFn == nil {
		return fmt.Errorf("unknown image registry type %s", c.RegistryType)
	}

	if c.imageReg, err = imageRegFn(c.Host(), c.httpClient); err != nil {
		return err
	}
	return nil
}

func (c *Config) initFilters() (err error) {
	if len(c.IncludeFilters) != 0 {
		if c.includeEngine, err = fl.New(c.IncludeFilters); err != nil {
			return err
		}
	}

	if len(c.ExcludeFilters) != 0 {
		if c.excludeEngine, err = fl.New(c.ExcludeFilters); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) initSkipList() (err error) {
	if c.SkipListPath != "" {
		if c.skipList, err = h.ReadLines(c.SkipListPath); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) initHTTPClient() (err error) {
	hcOptions := http.Option{AllowInsecureSSL: c.AllowInsecure}
	// prioritizing service path setups
	if c.ServiceAccountPath != "" {
		tokenSource := reg.TokenSourceMapper[c.RegistryType]
		if tokenSource == nil {
			return fmt.Errorf("cannot use %s for service account method", c.RegistryType)
		}
		data, err := ioutil.ReadFile(c.ServiceAccountPath)
		if err != nil {
			return err
		}
		ts, err := tokenSource(data)
		if err != nil {
			return err
		}
		hcOptions.TokenSource = ts
		if c.httpClient, err = http.New(hcOptions); err != nil {
			return err
		}
	} else {
		if c.RegUsername != "" && c.RegPassword == "" {
			return fmt.Errorf("you must set registry password")
		}

		if c.RegPassword != "" && c.RegUsername == "" {
			return fmt.Errorf("you must set registry username")
		}

		hcOptions.BasicAuth = struct {
			Username string
			Password string
		}{c.Username(), c.Password()}
		if c.httpClient, err = http.New(hcOptions); err != nil {
			return err
		}
	}
	return err
}

func (c *Config) initRepositoryList() (err error) {
	if c.RepoListPath == "" {
		return nil
	}
	data, err := ioutil.ReadFile(c.RepoListPath)
	if err != nil {
		return fmt.Errorf("error while reading repository list file: %w", err)
	}
	if err = json.Unmarshal(data, &c.repositories); err != nil {
		return fmt.Errorf("unmarshaling repository list file: %w", err)
	}
	return nil
}
