package server

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/rs/cors"
)

type Server *iris.Application

func GetOne(ctx context.Context, list []Summary) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StatusCode(400)
	}
	ctx.JSON(list[id-1])
}

func Start(port int) Server {
	app := iris.Default()

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}

	corsWrapper := cors.New(corsOptions).ServeHTTP

	app.WrapRouter(corsWrapper)

	app.Get("/", func(ctx context.Context) {
		ctx.StatusCode(200)
	})

	app.Get("/api/characters", func(ctx context.Context) {
		ctx.JSON(context.Map{"characters": Characters})
	})

	app.Get("/api/characters/:id", func(ctx context.Context) {
		GetOne(ctx, Characters)
	})

	app.Get("/api/planets", func(ctx context.Context) {
		ctx.JSON(context.Map{"planets": Planets})
	})

	app.Get("/api/planets/:id", func(ctx context.Context) {
		GetOne(ctx, Planets)
	})

	app.Get("/api/starships", func(ctx context.Context) {
		ctx.JSON(context.Map{"starships": Starships})
	})

	app.Get("/api/starships/:id", func(ctx context.Context) {
		GetOne(ctx, Starships)
	})

	app.Run(iris.Addr(":8080"))

	return app
}
