package main

import (
	"fmt"
	"log"
	"net/rpc"
)


type Move struct {
	Color int
	Col   int
}

type Board struct {
	BoardString string
}

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	var move Move

	fmt.Print("Enter color (0 = white, 1 = black): ")
	fmt.Scan(&move.Color)

	fmt.Print("Enter column: ")
	fmt.Scan(&move.Col)

	var replyMove int
	err = client.Call("ConnectGame.Move", move, &replyMove)
	if err != nil {
		log.Fatal("RPC error:", err)
	}

	log.Println("Sent Move RPC")

    var replyGet Board
	var args int
    err = client.Call("ConnectGame.Get", args, &replyGet)
    if err != nil {
        log.Fatal("game error:", err)
    }
    fmt.Printf("Game: \n%v", replyGet.BoardString)
}