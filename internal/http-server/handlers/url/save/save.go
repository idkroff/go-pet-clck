package save

import (
	"errors"
	resp "go-clck/internal/lib/api/response"
	"go-clck/internal/lib/logger/sl"
	"go-clck/internal/lib/random"
	"go-clck/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(savingUrl string, alias string) (int64, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
