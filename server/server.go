package server

import (
	"strings"

	"github.com/bryan-laipple/star-wars-server/storage"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/rs/cors"
)

type Server *iris.Application

func Start(port string) Server {

	swDbClient := storage.NewStarWarsDBClient()

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
		ctx.JSON(context.Map{"characters": swDbClient.GetCharacters()})
	})

	app.Get("/api/characters/:id", func(ctx context.Context) {
		id := ctx.Params().Get("id")
		if character, ok := swDbClient.GetCharacter(id); ok {
			ctx.JSON(character)
		} else {
			ctx.StatusCode(404)
		}
	})

	app.Get("/api/planets", func(ctx context.Context) {
		ctx.JSON(context.Map{"planets": swDbClient.GetPlanets()})
	})

	app.Get("/api/planets/:id", func(ctx context.Context) {
		id := ctx.Params().Get("id")
		if planet, ok := swDbClient.GetPlanet(id); ok {
			ctx.JSON(planet)
		} else {
			ctx.StatusCode(404)
		}
	})

	app.Get("/api/starships", func(ctx context.Context) {
		ctx.JSON(context.Map{"starships": swDbClient.GetStarships()})
	})

	app.Get("/api/starships/:id", func(ctx context.Context) {
		id := ctx.Params().Get("id")
		if starship, ok := swDbClient.GetStarship(id); ok {
			ctx.JSON(starship)
		} else {
			ctx.StatusCode(404)
		}
	})

	app.Run(iris.Addr(strings.Join([]string{"", port}, ":")))

	return app
}
