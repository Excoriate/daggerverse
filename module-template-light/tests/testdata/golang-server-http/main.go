package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Product represents a fake product struct
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

// fetchProducts handles the first GET API that returns a fake struct of products
func fetchProducts(c echo.Context) error {
	products := []Product{
		{ID: 1, Name: "Product A", Price: "$10"},
		{ID: 2, Name: "Product B", Price: "$20"},
		{ID: 3, Name: "Product C", Price: "$30"},
	}
	return c.JSON(http.StatusOK, products)
}

// fetchComments handles the second GET API that calls the public API
func fetchComments(c echo.Context) error {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to call external API")
	}
	defer resp.Body.Close()

	var comments interface{}
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse response from external API")
	}
	return c.JSON(http.StatusOK, comments)
}

// echoFlagValue handles the endpoint that returns the flag value
func echoFlagValue(c echo.Context) error {
	flagValue := c.Get("flagValue").(string)
	return c.String(http.StatusOK, flagValue)
}

func main() {
	// Define and parse the flag
	var flagValue string
	flag.StringVar(&flagValue, "flag", "default_value", "a string flag")
	flag.Parse()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Middleware to set flag value in context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("flagValue", flagValue)
			return next(c)
		}
	})

	e.GET("/products", fetchProducts)
	e.GET("/comments", fetchComments)
	e.GET("/flag", echoFlagValue)

	// Listen on all network interfaces
	address := "0.0.0.0:8080"
	e.Logger.Fatal(e.Start(address))
}
