package main

import "github.com/krakn/expense-management-backend-go/api"

func main() {
	app := app.NewApplicationServer()
	app.Run()
}
