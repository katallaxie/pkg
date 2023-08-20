package main

import (
	"fmt"
	"log"
	"os"

	"github.com/katallaxie/pkg/ulid"
)

func main() {
	id, err := ulid.New()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", id)
}
