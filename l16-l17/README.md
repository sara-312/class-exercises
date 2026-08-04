# Transactions and Two-Phase Commit

In this exercise you will be implementing basic 2PC for atomic transactions.

The operations in your transactions will be apart of a banking system that allows users to add and subtract money from accounts.

You will be able to run transactions that resemble:
```
op T1 alice 100
op T1 bob -100
commit T1
```
This transaction has two operations: adding 100 to alice's account and subtracting 100 from bob's.

You want these transactions to be atomic.
You also want to ensure an account can never go into a negative balance!

## txn.go

Create a new go file called `txn.go` and paste the following code into this file.

This go code implements an "almost" working key-value store with support for transactions.

Operations are in the form:
```golang
type Operation struct {
	Key string
	Value int
}
```

For example, the operation `{Key:"Alice", Value:100}` should **ADD** 100 to alice's account balance.

The operation `{Key:"Alice", Value:-100}` should **SUBTRACT** 100 to alice's account balance.


```golang
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Args struct {
	Op Operation
	TxnID string
}

type Reply struct {
	Success bool
}

type Operation struct {
	Key string
	Value int
}

type TxnServer struct {
	store map[string]int
	temp map[string]int
	txn map[string][]Operation
	lock sync.Mutex
}

func (t *TxnServer) Prepare(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	
	log.Printf("Prepare to commit TXN %v, Operations: %v\n", args.TxnID,t.txn[args.TxnID])

	vote := true
    temp := make(map[string]int, len(t.store))
    for key, value := range t.store {
        temp[key] = value
    }

	for _,op := range(t.txn[args.TxnID]) {
		temp[op.Key] += op.Value
		
	}

    // TODO: populate reply.Success with the vote

	return nil
}

func (t *TxnServer) Abort(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.temp = nil
	log.Println("ABORTED", t.store)
	return nil
}

func (t *TxnServer) Commit(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.store = t.temp
	log.Println("COMMITTED", t.store)
	return nil
}

func (t *TxnServer) Operation(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	_, exists := t.txn[args.TxnID]
	if !exists {
		t.txn[args.TxnID] = make([]Operation,0)
	}

	t.txn[args.TxnID] = append(t.txn[args.TxnID], Operation{args.Op.Key,args.Op.Value})
	log.Printf("Received TXN ID: %v Op: %v\n",args.TxnID, args.Op)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run txn.go <port>")
	}

	server := &TxnServer{
		store: make(map[string]int),
		txn: make(map[string][]Operation),
	}

	rpc.Register(server)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		log.Fatal("listen error:", err)
	}

	fmt.Println("TXN server listening on port", os.Args[1])

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("server error:", err)
	}
}
```

> [!IMPORTANT] 
> Add logic to `Prepare` to vote no if this transaction would go into a negative balance.

## Client.go

The following code implements a transaction client/coordinator. 

Once you run the client you will be able to use the following commands:

```
  op <txnid> <key> <value>
  commit <txnid>
  quit
```

Through the command line you can submit operations with a specified **transaction ID**, account **key**, and **value**.

Then, you can try to **commit** a specified transaction ID which will try to apply all or nothing of the operations within the specified transaction.

```golang

package main

import (
	"log"
	"net/rpc"
	"fmt"
	"strings"
	"bufio"
	"os"
	"strconv"
)

type Args struct {
	Op Operation
	TxnID string
}

type Reply struct {
	Success bool
}

type Operation struct {
	Key string
	Value int
}

func SubmitOp(client *rpc.Client, op Operation, txnId string) {
	args := &Args{Op: op, TxnID: txnId}
	var reply Reply
	err := client.Call("TxnServer.Operation",args, &reply)
	if(err != nil) {
		log.Fatal(err)
	}
}

func TryCommit(clients []*rpc.Client, txnId string) {
    //TODO: complete sending RPCS to clients
}

func cli(clients []*rpc.Client) {
	fmt.Println("Commands:")
	fmt.Println("  op <txnid> <key> <value>")
	fmt.Println("  commit <txnid>")
	fmt.Println("  quit")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {

		case "op":
			if len(fields) != 4 {
				fmt.Println("Usage: op <txnid> <key> <value>")
				continue
			}

			value, err := strconv.Atoi(fields[3])
			if err != nil {
				fmt.Println("value must be an integer")
				continue
			}

			op := Operation{
				Key:   fields[2],
				Value: value,
			}

			SubmitOp(clients[int(op.Key[0]) % len(clients)], op, fields[1])
			fmt.Printf("submitted to client %v\n", int(op.Key[0]) % len(clients))

		case "commit":
			if len(fields) != 2 {
				fmt.Println("Usage: commit <txnid>")
				continue
			}

			TryCommit(clients, fields[1])

		case "quit", "exit":
			fmt.Println("bye")
			return

		default:
			fmt.Println("Unknown command")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	
	txn_servers := []string{"localhost:1234","localhost:1235"}

	clients := make([]*rpc.Client,0)

	for _,address := range(txn_servers) {
		client, err := rpc.DialHTTP("tcp", address)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		clients = append(clients, client)
	}
	
	cli(clients)
}

```

> [!IMPORTANT] 
> Complete the implementation of `TryCommit` that runs Two Phase Commit across all servers in the clients slice.
> - Make sure to log the result of each server's vote
> - Make sure to log the Abort/Commit decision made for the transaction

Note: Your code can assume all transaction servers are up and running and does not need to handle timeouts.

# Testing your Code

Now, let's test your code. You will need to open multiple terminal windows to run the transaction servers and a client.

## Test 1
First, run a single transaction server and a single client.

Make sure to correctly update `txn_servers` in `client.go`.

```
go run txn.go 1234
go run client.go
```

Observe that a transaction drawing a single account into the negative is aborted and a transaction that goes into the positive is successful:

```
$> go run client.go
Commands:
  op <txnid> <key> <value>
  commit <txnid>
  quit
> op A alice -100
submitted to client 0
> commit A
2026/08/03 13:04:00 0: VOTED NO!
2026/08/03 13:04:00 TXN completed
```

On the txn server you should observe:
```
$> go run txn.go 1234
TXN server listening on port 1234
2026/08/03 13:03:58 Received TXN ID: A Op: {alice -100}
2026/08/03 13:04:00 Prepare to commit TXN A, Operations: [{alice -100}]
2026/08/03 13:04:00 TXN A: vote NO, negative balance @ alice
2026/08/03 13:04:00 ABORTED map[]
```

## Test 2
First, now run two transaction servers and a single client.

Make sure to correctly update `txn_servers` in `client.go`.

```
go run txn.go 1234
go run txn.go 1235
go run client.go
```

Observe that a transaction where one account would go into the negative, and the other would stay positive is aborted on both transaction servers.

Your output on the client should resemble: 
```
$> go run client.go
Commands:
  op <txnid> <key> <value>
  commit <txnid>
  quit
> op A alice 100
submitted to client 1
> op A bob 100
submitted to client 0
> commit A
2026/08/03 13:07:33 0: VOTED YES
2026/08/03 13:07:33 1: VOTED YES
2026/08/03 13:07:33 TXN completed
> op B alice -200
submitted to client 1
> op B bob 200
submitted to client 0
> commit B
2026/08/03 13:07:42 0: VOTED YES
2026/08/03 13:07:42 1: VOTED NO!
2026/08/03 13:07:42 TXN completed
> op C alice -100
submitted to client 1
> op C bob 100
submitted to client 0
> commit C
2026/08/03 13:08:01 0: VOTED YES
2026/08/03 13:08:01 1: VOTED YES
2026/08/03 13:08:01 TXN completed
```

Your output on the first txn server should resemble: 
```
go run txn.go 1234
TXN server listening on port 1234
2026/08/03 13:07:31 Received TXN ID: A Op: {bob 100}
2026/08/03 13:07:33 Prepare to commit TXN A, Operations: [{bob 100}]
2026/08/03 13:07:33 COMMITTED map[bob:100]
2026/08/03 13:07:39 Received TXN ID: B Op: {bob 200}
2026/08/03 13:07:42 Prepare to commit TXN B, Operations: [{bob 200}]
2026/08/03 13:07:42 ABORTED map[bob:100]
2026/08/03 13:08:00 Received TXN ID: C Op: {bob 100}
2026/08/03 13:08:01 Prepare to commit TXN C, Operations: [{bob 100}]
2026/08/03 13:08:01 COMMITTED map[bob:200]
```

Your output on the second txn server should resemble: 
```
go run txn.go 1235
TXN server listening on port 1235
2026/08/03 13:07:29 Received TXN ID: A Op: {alice 100}
2026/08/03 13:07:33 Prepare to commit TXN A, Operations: [{alice 100}]
2026/08/03 13:07:33 COMMITTED map[alice:100]
2026/08/03 13:07:36 Received TXN ID: B Op: {alice -200}
2026/08/03 13:07:42 Prepare to commit TXN B, Operations: [{alice -200}]
2026/08/03 13:07:42 TXN B: vote NO, negative balance @ alice
2026/08/03 13:07:42 ABORTED map[alice:100]
2026/08/03 13:07:57 Received TXN ID: C Op: {alice -100}
2026/08/03 13:08:01 Prepare to commit TXN C, Operations: [{alice -100}]
2026/08/03 13:08:01 COMMITTED map[alice:0]
```


# Warnings and Future Improvements

Consider: Is this transaction server serializable? 

This transaction system is bare-bones and can not properly handle concurrent transactions...

How would you incorporate 2PL into your `txn.go` Prepare code such that it implements serializable transactions?

Additionally, this client/coordinator does not handle the case where a transaction server is unreachable. Consider: How would you implement a timeout during the prepare phase? During the commit/abort phase?

