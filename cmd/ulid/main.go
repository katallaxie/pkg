package main

import (
	"fmt"
	"log"
	"os"

	"github.com/katallaxie/pkg/ulid"
)

func main() {
	id, err := ulid.NewReverse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", id)
}
