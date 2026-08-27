package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	domainErrors "github.com/KotelCatherine/mini-avito/internal/domain/errors"
	"github.com/KotelCatherine/mini-avito/internal/service"
	"github.com/go-chi/chi/v5"
)

type AdsHandler struct {
	service service.AdsService
}

func NewAdsHandler(serv service.AdsService) *AdsHandler {
	return &AdsHandler{service: serv}
}

func (adHandler *AdsHandler) CreatedAd(respWriter http.ResponseWriter, req *http.Request) {

	var createReq service.CreateAdRequest
	if err := json.NewDecoder(req.Body).Decode(&req); err != nil {
		http.Error(respWriter, "invalid request body", http.StatusBadRequest)
		return
	}

	ad, err := adHandler.service.CreateAd(req.Context(), createReq)
	if err != nil {

		switch {
		case errors.Is(err, domainErrors.ErrInvalidInput),
			errors.Is(err, domainErrors.ErrInvalidTitle),
			errors.Is(err, domainErrors.ErrAdInvalidPrice),
			errors.Is(err, domainErrors.ErrInvalidCategory):
			http.Error(respWriter, err.Error(), http.StatusBadRequest)
		default:
			http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		}

		return

	}

	respWriter.Header().Set("Content-Type", "applycation/json")
	respWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(respWriter).Encode(ad)

}

func (adHandler *AdsHandler) GetAd(respWriter http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")

	ad, err := adHandler.service.GetAd(req.Context(), id)
	if err != nil {

		switch {
		case errors.Is(err, domainErrors.ErrInvalidInput):
			http.Error(respWriter, "invalid id", http.StatusBadRequest)
		case errors.Is(err, domainErrors.ErrAdNotFound):
			http.Error(respWriter, "ad not found", http.StatusNotFound)
		default:
			http.Error(respWriter, "internal server error", http.StatusInternalServerError)
		}

		return

	}

	respWriter.Header().Set("Content-Type", "applycation/json")
	respWriter.WriteHeader(http.StatusAccepted)
	json.NewEncoder(respWriter).Encode(ad)

}
