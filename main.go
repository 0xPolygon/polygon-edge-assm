package main

import (
	lambda2 "Trapesys/polygon-edge-assm/framework/adapters/left/lambda"
	"Trapesys/polygon-edge-assm/framework/adapters/right/localstorage"
	app2 "Trapesys/polygon-edge-assm/internal/adapters/app"
	core2 "Trapesys/polygon-edge-assm/internal/adapters/core"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	core := core2.NewAdapter()
	lambdaAPI := lambda2.NewAdapter()
	localStorage := localstorage.NewAdapter()
	app := app2.NewAdapter(core, lambdaAPI, localStorage)
	lambda.Start(app.LambdaHandler)
}
