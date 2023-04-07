package main

import (
	"generadorSGA-go/internal/tags/handler"
	"generadorSGA-go/internal/tags/repository"
	"generadorSGA-go/internal/tags/service"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	r := repository.New()
	s := service.New(&r)
	h := handler.New(&s)
	lambda.Start(h.GetTags)
}
