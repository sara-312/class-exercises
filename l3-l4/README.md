# Week 2: RPC Server + Data Races

This week we will be making a simple "Connect 4" multiplayer game.

You will be programming a game **server** as well as a **client** program.

The server will:
- Take client requests for MOVE
    - A MOVE takes as input a COL and COLOR to drop the tile into.
    - The server verifies it is the correct color's TURN.
    - The server implements the logic of placing the piece in the correct place on the board.
- Take client requests for GET
    - GET returns the status of the current board

The client will:
- Allow users to chose a color
- Take user inputs from the command line
- Send MOVE requests to the server
- Send GET requests to the server and print the board to STDOUT for the user.

## Go net/rpc package

We will be using Go's `net/rpc` package to implement the networked communication.

https://pkg.go.dev/net/rpc 

## RPC Server

Copy the following file into `server.go`.

```golang
package main

import (
    "log"
	"net"
    "net/rpc"
	"net/http"
)

type Move struct {
	Color int
    Col int
}

type Board struct {
	BoardString string
}

type ConnectGame int

func (t *ConnectGame) Move(args *Move, reply *int) error {
	return nil
}

func (t *ConnectGame) Get(args *int, reply *Board) error {
    reply.BoardString = "Hello World"
	return nil
}

func main() {
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
```

## RPC Client

Copy the following starter code into a file called `client.go`.

```golang

package main

import (
    "log"
    "net/rpc"
)

type Move struct {
	Color int
    Col int
}

type Board struct {
	BoardString string
}

func main() {
    client, err := rpc.DialHTTP("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    // Synchronous call
    var reply Board
	var args int
    err = client.Call("ConnectGame.Get", &args, &reply)
    if err != nil {
        log.Fatal("game error:", err)
    }
   	log.Printf("Game: %v", reply)
}
```

Now, you should be able to run the client in one terminal window and the server in another terminal window via the following commands:

`go run server.go`

`go run client.go`

You should see the following output from the client:

```
2026/06/24 11:36:27 Game: {Hello World}
```

## Implementing Game

Now, let's implement the logic of the game server for `Move` and `Get`.

> [!IMPORTANT]
> 1. Chose a data structure to represent the board.
> 2. When the client places a piece in a column, implement the logic of checking the smallest occupied.

**Hint:** A good data structure for the board may be a global variable, 2D slice. You can set up a 2D slice as follows:

```go
var gameBoard [][]int
//...

func main() {
    gameBoard = make([][]int, rows) // Allocates the outer slice

	for i := range gameBoard {
		gameBoard[i] = make([]int, cols) // Allocates each inner row
	}
}
```

**Do not add any "win condition" checks.** If you finish all the sections of the README and want to implement all the connect4 logic, go ahead! :)

Here is some updated `client.go` code so that you can test your server:

```golang
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
```

In one tab, run:

```
go run server.go
```

In another terminal tab, each time you want to make a move run:

```
go run client.go
```

## Testing with Concurrent Clients

Now, what happens if we have several clients making moves at the same time? Will anything go wrong?

> [!IMPORTANT]
> Let's implement the following rule in `server.go`:
> - It should be **impossible** to have the same color move twice in a row.
>   - (i.e. if BLACK was the last color to move, we should only accept MOVE requests for the color WHITE).
> - If a move comes in that is for the wrong color, **return an error**.
 
Hint: You can use the `errors` package to throw custom errors in your RPC function body in `server.go`:
```return errors.New("Turn Order Violation")```

Now, let's implement the following test in `concurrent.go` to check that double-same-color moves are prohibited:

```golang
package main

import (
	"log"
	"net/rpc"
	"time"
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

    for i := range(10) {
        go func() {
            var reply int
            moveWhite := Move{0,i % 5}
            errWhite := client.Call("ConnectGame.Move", &moveWhite, &reply)
            if errWhite != nil {
                log.Println("RPC error:", errWhite)
            }
        }();
    }
    var replyB int
    moveBlack := Move{1,0}
    errBlack := client.Call("ConnectGame.Move", &moveBlack, &replyB)
    if errBlack != nil {
        log.Println("RPC error:", errBlack)
    }
	time.Sleep(10*time.Second)
}
```

You should see some:
`RPC error:` messages appear.

However, we want to be sure that our test works exactly as we expect.

**Does the server receive 10 requests for white moves and then 1 request for a black move?**

> [!IMPORTANT] 
> Add `log.Print` statements to inspect the order that Move requests are sent by the client and received by the server. Add a `log.Println` statement each time the server receives a `Move` RPC and print out the arguments of that RPC. 

**You should see that the Move request for black does not always come last!**

This is because when we call `go func()` 10 times, we just **start** these goroutines. We do not **wait** for any of the goroutines to run or finish their `Call`.

Thus, to fix our test, we need to use concurrency control!

## Concurrency Control: Channels

Let's use *channels* to ensure that we get 10 responses from our other threads before we send a move for BLACK.

Change `concurrent.go` to the following code. Fill in the correct lines to send/receive on channel `ch`. 

```golang
func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	ch := make(chan int)
    for i := range(10) {
        go func() {
            var reply int
            moveWhite := Move{0,i % 5}
            errWhite := client.Call("ConnectGame.Move", &moveWhite, &reply)
            if errWhite != nil {
                log.Println("RPC error:", errWhite)
				// send on channel
            } else {
				// send on channel
			}
        }();
    }
	sum := 0
	for range(10) {
		// receive on channel and update sum correctly to count successful moves
	}
	log.Println("Successful Moves:", sum)
    var replyB int
    moveBlack := Move{1,0}
    errBlack := client.Call("ConnectGame.Move", &moveBlack, &replyB)
    if errBlack != nil {
        log.Println("RPC error:", errBlack)
    }
}
```

Now, re-run your code with log.Print statements and check, does the BLACK move always appear last?

You should see that yes, it does always come last.

## Data Races

Make sure to run your server.go with `-race` flag.

`go run -race server.go`

For the purpose of this test, we will artificially encourage the data race to occur by adding a small wait before we update the game state recording who moved last.

> [!IMPORTANT]
> Add a `time.Sleep(10 * time.Millisecond)` call right before setting your `lastColorMoved` variable.

Is more than one successful moves ever reported?

You should see that sometimes our server appears to "glitch" and allow **many** white moves in a row.

What is happening?

If you are running your server with -race flag you should see output like:

```
WARNING: DATA RACE
Read at 0x00c0001de410 by goroutine 10:
  main.(*ConnectGame).Move()
      /Users/annaad/Documents/Work/bu/cs351/cs351-summer-26-mini/day2/server.go:49 +0x360
  runtime.call32()
      /usr/local/go/src/runtime/asm_arm64.s:670 +0x6c
  reflect.Value.Call()
      /usr/local/go/src/reflect/value.go:369 +0x90
  net/rpc.(*service).call()
      /usr/local/go/src/net/rpc/server.go:383 +0x1d8
  net/rpc.(*Server).ServeCodec.gowrap1()
      /usr/local/go/src/net/rpc/server.go:480 +0xb0

Previous write at 0x00c0001de410 by goroutine 15:
  main.(*ConnectGame).Move()
      /Users/annaad/Documents/Work/bu/cs351/cs351-summer-26-mini/day2/server.go:57 +0x488
  runtime.call32()
      /usr/local/go/src/runtime/asm_arm64.s:670 +0x6c
  reflect.Value.Call()
      /usr/local/go/src/reflect/value.go:369 +0x90
  net/rpc.(*service).call()
      /usr/local/go/src/net/rpc/server.go:383 +0x1d8
  net/rpc.(*Server).ServeCodec.gowrap1()
      /usr/local/go/src/net/rpc/server.go:480 +0xb0

Goroutine 10 (running) created at:
  net/rpc.(*Server).ServeCodec()
      /usr/local/go/src/net/rpc/server.go:480 +0x474
  net/rpc.(*Server).ServeConn()
      /usr/local/go/src/net/rpc/server.go:455 +0x570
  net/rpc.(*Server).ServeHTTP()
      /usr/local/go/src/net/rpc/server.go:710 +0x4f8
  net/http.(*ServeMux).ServeHTTP()
      /usr/local/go/src/net/http/server.go:2828 +0x1a8
  net/http.serverHandler.ServeHTTP()
      /usr/local/go/src/net/http/server.go:3311 +0x268
  net/http.(*conn).serve()
      /usr/local/go/src/net/http/server.go:2073 +0x9b4
  net/http.(*Server).Serve.gowrap3()
      /usr/local/go/src/net/http/server.go:3464 +0x40

Goroutine 15 (running) created at:
  net/rpc.(*Server).ServeCodec()
      /usr/local/go/src/net/rpc/server.go:480 +0x474
  net/rpc.(*Server).ServeConn()
      /usr/local/go/src/net/rpc/server.go:455 +0x570
  net/rpc.(*Server).ServeHTTP()
      /usr/local/go/src/net/rpc/server.go:710 +0x4f8
  net/http.(*ServeMux).ServeHTTP()
      /usr/local/go/src/net/http/server.go:2828 +0x1a8
  net/http.serverHandler.ServeHTTP()
      /usr/local/go/src/net/http/server.go:3311 +0x268
  net/http.(*conn).serve()
      /usr/local/go/src/net/http/server.go:2073 +0x9b4
  net/http.(*Server).Serve.gowrap3()
      /usr/local/go/src/net/http/server.go:3464 +0x40
==================
```

Notice it says we have a concurrent read and write in Move() at line 49 and 57. Check your own code for hints on which lines contain a concurrent read/write of the same data.

What we have found is an example of a data race!

Where are the two concurrently modifying threads? Each RPC is it's own goroutine in net/rpc package.

Thus, if two clients call Move, these Move RPCs will run concurrently on the server.

So if you have any code that looks like:
```golang
if lastColorMoved != args.Color {
    lastColorMoved = args.Color;
    //...   
}
```
It is possible that we have a **race** where both threads pass the if condition check before anyone updates `lastColorMoved`.

Thus, we need to add **concurrency control** to our RPC endpoints in our server code!

## Concurrency Control: Mutexes

You last task is to add a mutex to server.go to prevent data races.

Both the RPCs for Move/Get should acquire/release a lock around their body code, disallowing concurrent reads and writes of the game state data. The game state data includes the board, the color that last moved, etc.

https://gobyexample.com/mutexes

Here is some useful syntax to reference:
```golang
var mu sync.Mutex

// ...

mu.Lock()
defer mu.Unlock()
```

Hint: Make sure to import the "sync" package!

Now, even with a call to `time.Sleep` only one successful move should ever be reported!

