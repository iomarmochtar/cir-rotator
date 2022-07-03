package usecases

import (
	"fmt"

	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/rs/zerolog/log"
)

func DeleteRepositories(r reg.ImageRegistry, repositories []reg.Repository, skipList []string, isDryRun bool) (err error) {
	for idr := range repositories {
		repo := repositories[idr]
		if image, skipped := isInSkipList(repo, skipList); skipped {
			log.Debug().Str("image", image).Msg("listed in skip list, ignoring")
			continue
		}
		log.Info().Str("repo", repo.Name).Msg("deleting repository")
		if err = r.Delete(repo, isDryRun); err != nil {
			return err
		}
	}
	return nil
}

func isInSkipList(repo reg.Repository, skipList []string) (matches string, skip bool) {
	if len(skipList) == 0 {
		return "", false
	}

	for idd := range repo.Digests {
		digest := repo.Digests[idd]
		for idt := range digest.Tag {
			imageName := fmt.Sprintf("%s:%s", repo.Name, digest.Tag[idt])
			if h.IsInList[string](imageName, skipList) {
				return imageName, true
			}
		}
	}

	return "", false
}
