package app

import (
	"context"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/http/handlers"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/config"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/noteoutbox"
)

type NoteApplication struct {
	Config   *config.Config
	Services struct {
		note note.Service
	}
	Adapters struct {
		HTTP     *http.Server
		Postgres *postgres.Pool
		Auth     auth.Client
	}
}

func NewAppNote(ctx context.Context) (app *NoteApplication, err error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger.SetupLogger(cfg.Common.Logger)
	pgPool, err := postgres.New(ctx, cfg.Adapters.Postgres)
	if err != nil {
		return nil, err
	}
	noteOutBoxRepo := noteoutbox.NewRepository(pgPool)
	noteRepo := repository.NewRepository(
		&repository.Config{CreateNotesBatchSize: cfg.Services.Note.CreateNotesBatchSize},
		pgPool, noteOutBoxRepo,
	)
	noteService := note.NewService(cfg.Services.Note, noteRepo)
	noteHandlers := handlers.New(noteService)
	authService := auth.NewWrapper(cfg.Adapters.Auth)
	httpServer := http.NewServer(cfg.Adapters.HTTP, authService)
	app = new(NoteApplication)
	app.Config = cfg
	app.Services.note = noteService
	app.Adapters.HTTP = httpServer
	app.Adapters.Postgres = pgPool
	app.Adapters.Auth = authService
	app.initHTTPEndpoints(noteHandlers)
	return app, nil
}

func (n *NoteApplication) initHTTPEndpoints(h *handlers.NoteHandlers) {
	endpoints := []http.Endpoint{
		{Method: "POST", Path: "/create", Auth: true, Handler: h.CreateNote},
		{Method: "GET", Path: "/read", Auth: true, Handler: h.ReadNoteByID},
		{Method: "POST", Path: "/update", Auth: true, Handler: h.UpdateNote},
		{Method: "POST", Path: "/delete", Auth: true, Handler: h.DeleteNoteByID},
		{Method: "GET", Path: "/search", Auth: true, Handler: h.SearchNote},
	}

	for _, e := range endpoints {
		n.Adapters.HTTP.AddEndpoint(e)
	}
}
