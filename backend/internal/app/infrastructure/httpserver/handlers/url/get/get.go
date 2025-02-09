package get

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mixturka/DockerLens/backend/internal/app/application/interfaces"
	"github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/httpserver/response"
	"github.com/go-chi/render"
)

func NewGetHandler(log *slog.Logger, repo interfaces.PingRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		pings, err := repo.GetAllPings(ctx)
		if err != nil {
			log.Error("Failed GET request ", slog.Any("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Err("Data access failure"))
			return
		}

		render.JSON(w, r, map[string]any{"pings": pings})
	}
}
