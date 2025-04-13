package main

import (
	"context"
	"io"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
}
