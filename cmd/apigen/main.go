package main

import (
	"net/http"
	"os"

	"github.com/shaharby7/playground/pkg/greeter"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

func main() {
	reflector := openapi31.Reflector{}
	reflector.Spec = &openapi31.Spec{Openapi: "3.1.3"}
	reflector.Spec.Info.
		WithTitle("My fancy API").
		WithVersion("1.2.3").
		WithDescription("Allows you to greet people fancily")
	greetOp, err := reflector.NewOperationContext(http.MethodPost, "/api/greet")
	if err != nil {
		panic(err.Error())
	}

	greetOp.AddReqStructure(new(greeter.GreetInput))
	greetOp.AddRespStructure(new(greeter.GreetOutput), func(cu *openapi.ContentUnit) { cu.HTTPStatus = http.StatusOK })

	reflector.AddOperation(greetOp)

	schema, err := reflector.Spec.MarshalJSON()
	if err != nil {
		panic(err.Error())
	}
	os.WriteFile("./api/openapi.json", schema, 0777)
}
