package api

import (
	"fmt"

	auth "github.com/acs-dl/auth-svc/middlewares"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/acs-dl/mail-module-svc/internal/data/postgres"
	"github.com/acs-dl/mail-module-svc/internal/service/api/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (r *Router) apiRouter() chi.Router {
	router := chi.NewRouter()

	logger := r.cfg.Log().WithField("service", fmt.Sprintf("%s-api", data.ModuleName))

	secret := r.cfg.JwtParams().Secret

	router.Use(
		ape.RecoverMiddleware(logger),
		ape.LoganMiddleware(logger),
		ape.CtxMiddleware(
			//base
			handlers.CtxLog(logger),

			// storage
			handlers.CtxPermissionsQ(postgres.NewPermissionsQ(r.cfg.DB())),
			handlers.CtxUsersQ(postgres.NewUsersQ(r.cfg.DB())),
			handlers.CtxLinksQ(postgres.NewLinksQ(r.cfg.DB())),

			// connectors

			// other configs
			handlers.CtxParentContext(r.ctx),
		),
	)

	router.Route("/integrations/mail", func(r chi.Router) {
		r.With(auth.Jwt(secret, data.ModuleName, []string{"write", "read"}...)).
			Get("/get_input", handlers.GetInputs)

		r.Route("/links", func(r chi.Router) {
			r.With(auth.Jwt(secret, data.ModuleName, []string{"write"}...)).
				Post("/", handlers.AddLink)
			r.With(auth.Jwt(secret, data.ModuleName, []string{"write"}...)).
				Delete("/", handlers.RemoveLink)
		})

		r.With(auth.Jwt(secret, data.ModuleName, []string{"write", "read"}...)).
			Get("/submodule", handlers.CheckSubmodule)

		r.Get("/roles", handlers.GetRolesMap)          // comes from orchestrator
		r.Get("/user_roles", handlers.GetUserRolesMap) // comes from orchestrator

		r.Route("/estimate_refresh", func(r chi.Router) {
			r.Post("/submodule", handlers.GetEstimatedRefreshSubmodule)
			r.Post("/module", handlers.GetEstimatedRefreshModule)
		})

		r.With(auth.Jwt(secret, data.ModuleName, []string{"write", "read"}...)).
			Get("/permissions", handlers.GetPermissions)

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", handlers.GetUserById) // comes from orchestrator

			r.With(auth.Jwt(secret, data.ModuleName, []string{"write", "read"}...)).
				Get("/", handlers.GetUsers)
		})
	})

	return router
}
