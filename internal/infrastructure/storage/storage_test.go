package storage_test

import (
	"math/rand/v2"
	"testing"

	scrapperapi "github.com/es-debug/backend-academy-2024-go-template/api/openapi/v1/scrapper_api"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage"
	"github.com/stretchr/testify/require"
)

func TestDataInsertion(t *testing.T) {
	var (
		url     = "https://github.com/golang/go"
		tags    = []string{"tag"}
		filters = []string{"key:value"}
	)

	const (
		chatIDTagsFilters = iota
		chatIDTags
		chatIDFilters
		chatID
	)

	tests := map[string]struct {
		chatID int64
		link   models.Link
	}{
		"link with tags and filters insertion": {
			chatID: chatIDTagsFilters,
			link:   models.NewLink(rand.Int64(), url, tags, filters), //nolint:gosec // Temporary solution.
		},
		"link with tags insertion": {
			chatID: chatIDTags,
			link:   models.NewLink(rand.Int64(), url, tags, []string{}), //nolint:gosec // Temporary solution.
		},
		"link with filters insertion": {
			chatID: chatIDFilters,
			link:   models.NewLink(rand.Int64(), url, []string{}, filters), //nolint:gosec // Temporary solution.
		},
		"link without tags and filters": {
			chatID: chatID,
			link:   models.NewLink(rand.Int64(), url, []string{}, []string{}), //nolint:gosec // Temporary solution.
		},
	}

	repository := storage.NewLinksInMemoryService()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := repository.AddChat(test.chatID); err != nil {
				require.FailNow(t, err.Error())
			}

			if err := repository.AddLink(test.chatID, test.link); err != nil {
				require.FailNow(t, err.Error())
			}

			links, err := repository.GetLinks(test.chatID)
			if err != nil {
				require.FailNow(t, err.Error())
			}

			for _, link := range links {
				if *test.link.Id == *link.Id {
					require.Equal(t, test.link, link)
				}
			}
		})
	}
}

func TestHappyPath(t *testing.T) {
	var (
		url     = "https://github.com/golang/go"
		tags    = []string{"tag"}
		filters = []string{"key:value"}
	)

	const (
		chatIDTagsFilters = iota
		chatIDTags
		chatIDFilters
		chatID
	)

	tests := map[string]struct {
		chatID int64
		link   models.Link
	}{
		"link with tags and filters": {
			chatID: chatIDTagsFilters,
			link:   models.NewLink(rand.Int64(), url, tags, filters), //nolint:gosec // Temporary solution.
		},
		"link with tags": {
			chatID: chatIDTags,
			link:   models.NewLink(rand.Int64(), url, tags, []string{}), //nolint:gosec // Temporary solution.
		},
		"link with filters": {
			chatID: chatIDFilters,
			link:   models.NewLink(rand.Int64(), url, []string{}, filters), //nolint:gosec // Temporary solution.
		},
		"link without tags and filters": {
			chatID: chatID,
			link:   models.NewLink(rand.Int64(), url, []string{}, []string{}), //nolint:gosec // Temporary solution.
		},
	}

	repositories := storage.NewLinksInMemoryService()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := repositories.AddChat(test.chatID); err != nil {
				require.FailNow(t, err.Error())
			}

			if err := repositories.AddLink(test.chatID, test.link); err != nil {
				require.FailNow(t, err.Error())
			}

			if err := repositories.DeleteLink(test.chatID, *test.link.Url); err != nil {
				require.FailNow(t, err.Error())
			}
		})
	}
}

func TestDoubleInsertion(t *testing.T) {
	var (
		url     = "https://github.com/golang/go"
		tags    = []string{"tag"}
		filters = []string{"key:value"}
	)

	const (
		chatIDTagsFilters = iota
		chatIDTags
		chatIDFilters
		chatID
	)

	tests := map[string]struct {
		chatID int64
		link   models.Link
		want   error
	}{
		"link with tags and filters": {
			chatID: chatIDTagsFilters,
			link:   models.NewLink(rand.Int64(), url, tags, filters), //nolint:gosec // Temporary solution.
			want:   scrapperapi.ErrLinkAlreadyExists,
		},
		"link with tags": {
			chatID: chatIDTags,
			link:   models.NewLink(rand.Int64(), url, tags, []string{}), //nolint:gosec // Temporary solution.
			want:   scrapperapi.ErrLinkAlreadyExists,
		},
		"link with filters": {
			chatID: chatIDFilters,
			link:   models.NewLink(rand.Int64(), url, []string{}, filters), //nolint:gosec // Temporary solution.
			want:   scrapperapi.ErrLinkAlreadyExists,
		},
		"link without tags and filters": {
			chatID: chatID,
			link:   models.NewLink(rand.Int64(), url, []string{}, []string{}), //nolint:gosec // Temporary solution.
			want:   scrapperapi.ErrLinkAlreadyExists,
		},
	}

	repositories := storage.NewLinksInMemoryService()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := repositories.AddChat(test.chatID); err != nil {
				require.FailNow(t, err.Error())
			}

			if err := repositories.AddLink(test.chatID, test.link); err != nil {
				require.Equal(t, err, scrapperapi.ErrLinkAlreadyExists)
			}
		})
	}
}
