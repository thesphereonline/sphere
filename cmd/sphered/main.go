// cmd/sphered/main.go
package main

import (
	"sphere/internal/api"
	"sphere/internal/core"
)

func main() {
	bc := core.NewBlockchain()
	api.StartServer(bc)
}
