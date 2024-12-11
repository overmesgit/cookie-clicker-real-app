package main

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"kafka-example/handlers"
	_ "kafka-example/migrations"
	"log"
	"os"
	"strings"
)

func main() {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(
		app, app.RootCmd, migratecmd.Config{
			Automigrate: isGoRun,
		},
	)

	app.OnServe().BindFunc(
		func(e *core.ServeEvent) error {
			err := handlers.InitUserLevels(app)
			if err != nil {
				return err
			}
			return e.Next()
		},
	)

	app.OnServe().BindFunc(
		func(se *core.ServeEvent) error {
			se.Router.GET("/{path...}", apis.Static(os.DirFS("./views"), false))
			se.Router.POST(
				"/signup", handlers.HandleSignUp(app),
			)

			se.Router.POST(
				"/click", handlers.HandleClickEvent(app),
			).Bind(apis.RequireAuth())
			se.Router.GET(
				"/level", handlers.HandleLevelsList(app),
			).Bind(apis.RequireAuth())
			se.Router.POST(
				"/level", handlers.HandleLevelEvent(app),
			).Bind(apis.RequireAuth())
			se.Router.GET(
				"/upgrade", handlers.HandleUpgradesList(app),
			).Bind(apis.RequireAuth())
			se.Router.POST(
				"/upgrade", handlers.HandleUpgradeEvent(app),
			).Bind(apis.RequireAuth())

			return se.Next()
		},
	)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
