package app

import (
	"fmt"

	c "github.com/iomarmochtar/cir-rotator/app/config"
	fl "github.com/iomarmochtar/cir-rotator/pkg/filter"
	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/rs/zerolog/log"
)

const VERSION = "0.1.2"

type App struct {
	config c.IConfig
}

func New(config c.IConfig) *App {
	return &App{config}
}

func (a App) ListRepositories() ([]reg.Repository, error) {
	return a.fetchAndFilterRepositories()
}

func (a App) DeleteRepositories(repositories []reg.Repository) (err error) {
	skipList := a.config.SkipList()
	for idr := range repositories {
		repo := repositories[idr]
		// filter the list of tags if skiplist provided, if it's matched then ignore the related digest for deletion
		if len(skipList) != 0 {
			filterRepositoryDigestBySkipList(&repo, skipList)
		}
		// if there is no such digests in repository so then nothing todo with it.
		if len(repo.Digests) == 0 {
			log.Warn().Str("repo", repo.Name).Msg("no digest found as for deleting in repository, skip it")
			continue
		}

		if a.config.IsDryRun() {
			log.Info().Str("repo", repo.Name).Msg("[DRY_RUN] attempting for deletion")
			continue
		}

		log.Info().Str("repo", repo.Name).Msg("deleting matched images in the repository")
		if err = a.config.ImageRegistry().Delete(repo); err != nil {
			return err
		}
	}
	return nil
}

func (a App) fetchAndFilterRepositories() ([]reg.Repository, error) {
	repositories, err := a.config.ImageRegistry().Catalog()
	if err != nil {
		return nil, err
	}

	includeFilter := a.config.IncludeEngine()
	excludeFilter := a.config.ExcludeEngine()
	// if there is no filters then just return as is
	if includeFilter == nil && excludeFilter == nil {
		return repositories, nil
	}

	return doFilter(repositories, includeFilter, excludeFilter)
}

// filterRepositories filter listing repositories and digest based in include and exclude filters
func doFilter(repositories []reg.Repository, includeFilter, excludeFilter fl.IFilterEngine) ([]reg.Repository, error) {
	//nolint:prealloc
	var result []reg.Repository

	for idr := range repositories {
		repo := repositories[idr]
		var resultDigest []reg.Digest
		for idd := range repo.Digests {
			digest := repo.Digests[idd]
			fields := fl.Fields{
				Repository: repo.Name,
				Digest:     digest.Name,
				ImageSize:  digest.ImageSizeBytes,
				Tags:       digest.Tag,
				CreatedAt:  digest.Created,
				UploadedAt: digest.Uploaded,
			}

			if includeFilter != nil {
				iResult, err := includeFilter.Process(fields)
				if err != nil {
					return nil, err
				}
				if !iResult {
					continue
				}
			}

			if excludeFilter != nil {
				eResult, err := excludeFilter.Process(fields)
				if err != nil {
					return nil, err
				}
				if eResult {
					continue
				}
			}
			resultDigest = append(resultDigest, digest)
		}

		if len(resultDigest) == 0 {
			continue
		}
		result = append(result, reg.Repository{Name: repo.Name, Digests: resultDigest})
	}
	return result, nil
}

func filterRepositoryDigestBySkipList(repo *reg.Repository, skipList []string) {
	tmpDigests := []reg.Digest{}
	for idd := range repo.Digests {
		includeDigest := true
		digest := repo.Digests[idd]
		for idt := range digest.Tag {
			imageName := fmt.Sprintf("%s:%s", repo.Name, digest.Tag[idt])
			if h.IsInList(imageName, skipList) {
				includeDigest = false
				log.Info().Str("image", imageName).Str("digest", digest.Name).Msg("listed in skip list, ignoring related digest")
				break
			}
		}
		if includeDigest {
			tmpDigests = append(tmpDigests, digest)
		}
	}
	repo.Digests = tmpDigests
}
