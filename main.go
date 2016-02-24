package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gombadi/lambda-reflector/apigateway"
)

func main() {

	// decode the last commandline arg which is the api-gateway input
	ar, err := apigateway.DecodeRequest(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error: unable to create new lambda event: %v\n", err)
	}

	switch {
	// some test cases to check the mappings is working as expected
	case ar.GetValue("query.ret") == "503":
		ar.ExitOnErr("regex503", "invalid input request")
	case ar.GetValue("query.ret") == "404":
		ar.ExitOnErr("notFound", "requested object not found")
	case ar.GetValue("query.redir") == "123":
		ar.Redirect302("https://www.google.com/")
	case ar.GetValue("query.redir") == "456":
		ar.Redirect302("https://www.golang.org/")
	case ar.GetValue("query.type") == "all":
		if b, err := ar.GetJSON(); err != nil {
			fmt.Printf(" ")
		} else {
			fmt.Printf("%s", b)
		}
	default:
		fmt.Printf("%s", ar.GetValue("sourceip"))
	}
	// normal end will cause wrapper to context.succeed and api-gateway to return 200
}

/*

 */
