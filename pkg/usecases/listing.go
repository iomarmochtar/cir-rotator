package usecases

import (
	fl "github.com/iomarmochtar/cir-rotator/pkg/filter"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
)

func ListRepositories(r reg.ImageRegistry, includeFilter, excludeFilter fl.IFilterEngine) ([]reg.Repository, error) {
	repositories, err := r.Catalog()
	if err != nil {
		return nil, err
	}

	repositories, err = doFilter(repositories, includeFilter, excludeFilter)
	if err != nil {
		return nil, err
	}

	return repositories, nil
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
