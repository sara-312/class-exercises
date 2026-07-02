package main

import (
    "log"
	"net"
    "net/rpc"
	"net/http"
	"fmt"
)

type Move struct {
	Color int
    Col int
}

type Board struct {
	BoardString string
}
const rows = 6

const cols = 7

var gameBoard [][]int

type ConnectGame int

func (t *ConnectGame) Move(args *Move, reply *int) error {
	col := args.Col
	for i := rows - 1; i >= 0; i-- {
		if gameBoard[i][col] == 0 {
			gameBoard[i][col] = args.Color + 1
			*reply = i
			return nil
		}
	}
	return fmt.Errorf("Column %d is full", col)
}

func (t *ConnectGame) Get(args *int, reply *Board) error {
	boardString := ""
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			boardString += fmt.Sprintf("%d ", gameBoard[i][j])
		}
		boardString += "\n"
	}
	reply.BoardString = boardString
	return nil
}

func main() {
	gameBoard = make([][]int, rows)
	for i := range gameBoard {
		gameBoard[i] = make([]int, cols)
	}
    cg := new(ConnectGame)
    rpc.Register(cg)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("listen error:", err)
    }
	log.Println("Serving on PORT 1234")
    http.Serve(l, nil)
}