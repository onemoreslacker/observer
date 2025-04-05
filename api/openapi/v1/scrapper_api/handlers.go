package scrapperapi

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/entities"
)

type API struct {
	links linksRepository
}

type linksRepository interface {
	AddChat(id int64) error
	DeleteChat(id int64) error
	AddLink(id int64, link entities.Link) error
	GetLinks(id int64) ([]entities.Link, error)
	DeleteLink(id int64, url string) (entities.Link, error)
}

func New(links linksRepository) *API {
	return &API{
		links: links,
	}
}

//nolint:revive,stylecheck // Generated code cannot be edited.
func (a *API) PostTgChatId(w http.ResponseWriter, _ *http.Request, id int64) {
	if err := a.links.AddChat(id); err != nil {
		respondWithError(w, http.StatusBadRequest,
			err.Error(), ErrInvalidBody.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, http.NoBody)
}

//nolint:revive,stylecheck // Generated code cannot be edited.
func (a *API) GetTgChatId(w http.ResponseWriter, _ *http.Request, id int64) {
	if _, err := a.links.GetLinks(id); err != nil {
		respondWithError(w, http.StatusBadRequest,
			err.Error(), ErrInvalidBody.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, http.NoBody)
}

//nolint:revive,stylecheck // Generated code cannot be edited.
func (a *API) DeleteTgChatId(w http.ResponseWriter, _ *http.Request, id int64) {
	if err := a.links.DeleteChat(id); err != nil {
		respondWithError(w, http.StatusBadRequest,
			err.Error(), ErrInvalidBody.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, http.NoBody)
}

func (a *API) PostLinks(w http.ResponseWriter, r *http.Request, params PostLinksParams) {
	id := params.TgChatId

	var model AddLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		respondWithError(w, http.StatusBadRequest,
			err.Error(), ErrInvalidBody.Error())
		return
	}

	if model.Link == nil || model.Tags == nil || model.Filters == nil {
		respondWithError(w, http.StatusBadRequest,
			ErrAddLinkInvalidLink.Error(), ErrInvalidBody.Error())
		return
	}

	url, tags, filters := *model.Link, *model.Tags, *model.Filters

	if !checkResourceAvailability(url) {
		respondWithError(w, http.StatusBadRequest,
			ErrAddLinkInvalidLink.Error(), ErrInvalidBody.Error())
		return
	}

	link := entities.NewLink(rand.Int64(), url, tags, filters) //nolint:gosec // Temporary solution

	if err := a.links.AddLink(id, link); err != nil {
		var status int

		if errors.Is(err, ErrLinkAlreadyExists) {
			status = http.StatusConflict
		} else {
			status = http.StatusBadRequest
		}

		respondWithError(w, status, err.Error(),
			ErrAddLinkFailed.Error())

		return
	}

	respondWithJSON(w, http.StatusOK, LinkResponse(link))
}

func (a *API) GetLinks(w http.ResponseWriter, _ *http.Request, params GetLinksParams) {
	id := params.TgChatId

	links, err := a.links.GetLinks(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			ErrGetLinksFailed.Error(), ErrInvalidBody.Error())
		return
	}

	linksResponse := make([]LinkResponse, len(links))
	for i := range len(links) {
		linksResponse[i] = LinkResponse(links[i])
	}

	sz := int32(len(linksResponse)) //nolint:gosec // Generated code cannot be edited.

	respondWithJSON(w, http.StatusOK, ListLinksResponse{
		Links: &linksResponse,
		Size:  &sz,
	})
}

func (a *API) DeleteLinks(w http.ResponseWriter, r *http.Request, params DeleteLinksParams) {
	id := params.TgChatId

	var model RemoveLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		respondWithError(w, http.StatusBadRequest,
			err.Error(), ErrInvalidBody.Error())
		return
	}

	if model.Link == nil {
		respondWithError(w, http.StatusBadRequest,
			ErrDeleteLinkInvalidLink.Error(), ErrInvalidBody.Error())
		return
	}

	link, err := a.links.DeleteLink(id, *model.Link)
	if err != nil {
		var status int

		if errors.Is(err, ErrLinkAlreadyExists) {
			status = http.StatusConflict
		} else {
			status = http.StatusBadRequest
		}

		respondWithError(w, status, err.Error(),
			ErrAddLinkFailed.Error())

		return
	}

	respondWithJSON(w, http.StatusOK, LinkResponse(link))
}
