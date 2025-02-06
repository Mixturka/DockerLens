package get

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Mixturka/DockerLens/backend/internal/app/application/interfaces"
	"github.com/Mixturka/DockerLens/backend/internal/app/infrastructure/httpserver/response"
	"github.com/go-chi/render"
)

func NewGetHandler(log *slog.Logger, repo interfaces.PingRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		cursor := r.URL.Query().Get("cursor")
		limitStr := r.URL.Query().Get("limit")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 30
		}

		pings, nextCursor, err := repo.GetPingsCursor(ctx, limit, cursor)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Err("Data access failure"))
			return
		}

		resp := map[string]interface{}{
			"data":        pings,
			"next_cursor": nextCursor,
		}

		render.JSON(w, r, resp)
	}
}
