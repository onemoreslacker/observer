package core

import (
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	botclient "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/clients/bot"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/clients/external"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scrapper struct {
	botClient  botclient.ClientInterface
	repository LinksService
	external   ExternalClient
	sch        gocron.Scheduler
}

type LinksService interface {
	GetChatIDs() ([]int64, error)
	GetLinks(int64) (links []models.Link, err error)
}

type ExternalClient interface {
	RetrieveStackOverflowUpdates(link string) ([]models.StackOverflowUpdate, error)
	RetrieveGitHubUpdates(link string) ([]models.GitHubUpdate, error)
}

func New(client botclient.ClientInterface, repository LinksService, sch gocron.Scheduler) (*Scrapper, error) {
	return &Scrapper{
		botClient:  client,
		repository: repository,
		external:   external.New(),
		sch:        sch,
	}, nil
}

func (s *Scrapper) Run() error {
	_, err := s.sch.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(10, 0, 0),
			),
		),
		gocron.NewTask(
			func() error {
				return s.scrapeUpdates()
			},
		),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					slog.Error(
						"job error",
						slog.String("msg", err.Error()),
						slog.String("job_id", jobID.String()),
						slog.String("job_name", jobName),
						slog.String("service", "scrapper"),
					)
				},
			),
		),
	)

	if err != nil {
		return err
	}

	s.sch.Start()

	return nil
}
