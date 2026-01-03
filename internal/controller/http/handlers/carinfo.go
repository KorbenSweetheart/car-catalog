package handlers

import (
	"log/slog"
	"net/http"
	"viewer/internal/domain"
)

type CarProvider interface {
	CarModel(int) (domain.Car, error)
	CarManufacturer(int) (domain.Manufacturer, error)
	CarCategory(int) (domain.Category, error)
}

func NewCarModel(log *slog.Logger, carProvider CarProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getcar.New"

		log = log.With(
			slog.String("op", op),
		)

		// input := dto.GetCarInput{
		// 	ID: chi.URLParam(r, "id"),
		// }

		// output, err := h.profileService.GetProfile(r.Context(), input)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)

		// 	return
		// }

	}

}
