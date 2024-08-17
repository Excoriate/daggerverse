package main

import (
	"encoding/json"
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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/products", fetchProducts)
	e.GET("/comments", fetchComments)

	// Listen on all network interfaces
	address := "0.0.0.0:8080"
	e.Logger.Fatal(e.Start(address))
}
