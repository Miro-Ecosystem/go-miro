package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Miro-Ecosystem/go-miro/miro"
)

func main() {
	c := miro.NewClient(os.Getenv("MIRO_ACCESS_KEY"))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	name := os.Args[1]
	desc := os.Args[2]

	req := &miro.CreateBoardRequest{
		Name:        name,
		Description: desc,
	}

	board1, err := c.Boards.Create(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	board2, err := c.Boards.Get(ctx, board1.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s board created", board2.Name)

	c.Boards.Delete(ctx, board2.ID)
	fmt.Printf("%s board deleted", board2.Name)
}
