package main

import (
	"errors"
	"strings"

	"github.com/olivebay/bugify"
)

func returnError() error {
	return errors.New("this is an error")
}

func main() {
	client := bugify.NewClient(os.Getenv("CHECKLY_API_KEY"), "olivebay/urlinfo2")

	err := returnError()
	if err != nil{
		client.Create("failed to return an error %v", err)
	}
}
