package save

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mixturka/DockerLens/backend/internal/app/application/dtos"
	"github.com/Mixturka/DockerLens/backend/internal/app/application/interfaces"
	"github.com/Mixturka/DockerLens/backend/internal/app/domain/entities"
	"github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/httpserver/response"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewSaveHandler(log *slog.Logger, repo interfaces.PingRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dtos.PingDTO

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			// log.Error("failed to decode request body", err.Error())
			render.JSON(w, r, response.Err("failed to decode request body"))
			return
		}

		log.Info("request body decoded successfully", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			// log.Error("request validation failed", err.Error())
			render.Status(r, http.StatusNotAcceptable)
			render.JSON(w, r, response.Err("request validation failed"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err = repo.Save(ctx, entities.Ping{
			ID:        uuid.NewString(),
			IP:        req.IP,
			IsSuccess: req.IsSuccess,
			Time:      req.Time,
			CreatedAt: time.Now(),
		})

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Err("internal save failed: "+err.Error()))
			return
		}

		render.JSON(w, r, response.Ok())
	}
}
