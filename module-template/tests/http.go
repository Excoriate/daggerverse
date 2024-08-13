package main

import (
	"context"
	"encoding/json"
	"strings"
)

// TestHTTPCurl tests an HTTP GET request using the curl command within an Alpine container.
//
// This function performs the following steps:
// 1. Sets up the Alpine container with necessary utilities to perform the curl operation.
// 2. Executes the curl command against the specified target URL and captures the output.
// 3. Verifies that the curl command produced non-empty output.
// 4. Checks for errors during the curl command execution.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If the curl command fails or produces an empty output, an error is returned.
func (m *Tests) TestHTTPCurl(ctx context.Context) error {
	targetURL := "https://fakestoreapiserver.reactbd.com/smart"

	// Set up the Alpine container with required utilities for HTTP operations.
	targetModule := dag.ModuleTemplate()
	targetModule = targetModule.
		WithUtilitiesInAlpineContainer().
		WithHttpcurl(targetURL)

	// Execute the curl command and capture the output.
	out, err := targetModule.Ctr().Stdout(ctx)

	// Check if the output is empty, indicating a potential issue.
	if out == "" {
		return Errorf("failed to inspect the curl output of the URL %s. Got empty output", targetURL)
	}

	// Check for any error during the curl command execution.
	if err != nil {
		return WrapErrorf(err, "failed to curl this URL %s", targetURL)
	}

	return nil
}

// ProductJSONApiTest represents the structure of the product
// information returned by the API.
//
// Fields:
// - ID: Unique identifier for each product.
// - Title: Name/title of the product.
// - IsNew: Boolean indicating if the product is new.
// - OldPrice: The previous price of the product, represented as a string.
// - Price: The current price of the product.
// - Description: A brief description of the product.
// - Category: The category to which the product belongs.
// - Image: URL to the product's image.
// - Rating: Rating of the product out of 5.
type ProductJSONApiTest struct {
	ID          int     `json:"_id"` //nolint:tagliatelle // ID is a unique identifier for each product.
	Title       string  `json:"title"`
	IsNew       bool    `json:"isNew"`
	OldPrice    string  `json:"oldPrice"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      int     `json:"rating"`
}

// TestHTTPDoJSONAPICall tests an HTTP GET request to fetch product information from a JSON API.
//
// This function performs the following steps:
// 1. Sends an HTTP GET request to the specified URL to fetch product information in JSON format.
// 2. Checks if the response is non-nil.
// 3. Reads the contents of the JSON response file.
// 4. Verifies that the content is not empty and that it does not contain an error message.
// 5. Unmarshals the JSON response into a slice of ProductJSONApiTest structs.
// 6. Ensures that the unmarshalling was successful and the response contains at least one product.
//
// Parameters:
// - ctx: The context to control the execution.
//
// Returns:
// - error: If any of the steps fail, an error is returned.
func (m *Tests) TestHTTPDoJSONAPICall(ctx context.Context) error {
	targetURL := "https://fakestoreapiserver.reactbd.com/products"

	targetModule := dag.ModuleTemplate()
	jsonFile := targetModule.DoJsonapicall("GET", targetURL)

	if jsonFile == nil {
		return Errorf("failed to get the JSON response from the URL %s", targetURL)
	}

	content, err := jsonFile.Contents(ctx)
	if err != nil {
		return WrapErrorf(err, "failed to get the contents of the file /response.json")
	}

	if content == "" {
		return Errorf("failed to get the contents of the file /response.json")
	}

	// Check if the response is an error message
	if strings.Contains(content, "Bad request") || strings.Contains(content, "BAD_REQUEST") {
		return Errorf("API returned an error: %s", content)
	}

	// Unmarshal the JSON content into a slice of ProductJSONApiTest structs
	var products []ProductJSONApiTest
	err = json.Unmarshal([]byte(content), &products)

	if err != nil {
		return WrapErrorf(err, "failed to unmarshal the JSON response. Raw content: %s", content)
	}

	// Ensure that the unmarshalling was successful and the response contains at least one product
	if len(products) == 0 {
		return Errorf("failed to unmarshal the JSON response or the response was empty")
	}

	return nil
}
