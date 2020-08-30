package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Miro-Ecosystem/go-miro/miro"
)

func main() {
	c := miro.NewClient(os.Getenv("MIRO_ACCESS_KEY"))
	board, err := c.Boards.Get(context.Background(), os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(board.Name)
}
