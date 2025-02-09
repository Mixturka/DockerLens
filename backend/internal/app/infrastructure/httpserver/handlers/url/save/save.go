package save

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mixturka/DockerLens/backend/internal/app/application/dtos"
	"github.com/Mixturka/DockerLens/backend/internal/app/application/interfaces"
	"github.com/Mixturka/DockerLens/backend/internal/app/domain/entities"
	"github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/httpserver/response"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func NewSaveHandler(log *slog.Logger, repo interfaces.PingRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []dtos.PingDTO

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to decode request body: %s", err))
			render.JSON(w, r, response.Err("failed to decode request body"))
			return
		}

		log.Debug("request body decoded successfully", slog.Any("request", req))

		validate := validator.New()
		for _, ping := range req {
			if err := validate.Struct(ping); err != nil {
				log.Error("Validation failed", slog.String("error", err.Error()), slog.Any("ping", ping))
				render.Status(r, http.StatusNotAcceptable)
				render.JSON(w, r, response.Err("request validation failed"))
				return
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		pings := []entities.Ping{}
		for _, ping := range req {
			pings = append(pings, entities.Ping{
				IP:          ping.IP,
				IsSuccess:   ping.IsSuccess,
				Duration:    ping.Duration,
				LastSuccess: time.Now(),
			})
		}

		err = repo.SaveBatch(ctx, pings)
		if err != nil {
			log.Error("failed to save batch pings", slog.String("error", err.Error()))
			render.JSON(w, r, response.Err("failed to save pings"))
			return
		}
		render.JSON(w, r, response.Ok())
	}
}
