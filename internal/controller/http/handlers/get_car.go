package handlers

import (
	"log/slog"
	"net/http"
	"viewer/internal/dto"
)

type CarProvider interface {
	GetCarModel(id int) (dto.GetCarOutput, error) // add car as a struct or DTO
}

func New(log *slog.Logger, CarProvider CarProvider) http.HandlerFunc {
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
