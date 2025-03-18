package main

import (
	"os"
	"product-golang-graphql/auth"
	"product-golang-graphql/graph"
	"product-golang-graphql/graph/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := fiber.New()
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://gofiber.io, http://localhost:3000, http://localhost:4000, http://192.168.70.17:3000",
		AllowHeaders: "Origin, Content-Type, Accept,  Access-Control-Allow-Origin, Authorization",
	}))

	

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})



	// srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	// srv.Use(extension.AutomaticPersistedQuery{
	// 	Cache: lru.New[string](100),
	// })

	
	
	graphQLHandler := fasthttpadaptor.NewFastHTTPHandler(srv)

	app.Post("/graphql", func(c *fiber.Ctx) error {
		c.Request().Header.Set("Content-Type", "application/json")
		graphQLHandler(c.Context())
		return nil
	})

	playgroundHandler := fasthttpadaptor.NewFastHTTPHandler(playground.Handler("GraphQL Playground", "/graphql"))
	playgroundAuthHandler := fasthttpadaptor.NewFastHTTPHandler(playground.Handler("GraphQL Playground", "/graphql/auth"))

	app.Get("/playground", func(c *fiber.Ctx) error {
		playgroundHandler(c.Context())
		return nil
	})
	app.Get("/playground/auth", func(c *fiber.Ctx) error {
		playgroundAuthHandler(c.Context())
		return nil
	})

	// Endpoint untuk GraphQL dengan autentikasi
	app.Post("/graphql/auth", auth.JWTMiddleware(), func(c *fiber.Ctx) error {

		c.Request().Header.Set("Content-Type", "application/json")
		// ctx := context.WithValue(c.Context(), "fiberContext", c)
		user, ok := c.Locals("user").(*model.User)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		reqCtx := c.Context() // This is already *fasthttp.RequestCtx

	// Inject the updated context (with user)
		reqCtx.SetUserValue("user", user)
		// Simpan context Fiber di GraphQL Context
		graphQLHandler(reqCtx)
		return nil
	})

	// Jalankan server
	app.Listen(":"+port)

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
}
