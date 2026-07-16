# Heartbeats

**This is a group activity! Get into groups of 4-5 students!**

You will be writing `server.go` code that sends and receives heartbeats with other students' servers. 

You will implement a **gossip-based protocol** to disseminate information.

## server.go


Copy and paste the following starter code into `server.go`.

```golang
package main

import (
	"sync"
	"log"
	"net/rpc"
	"net"
	"net/http"
	"fmt"
)

type PeerServer struct {
	Address string
	Client *rpc.Client
}

type Args struct {
	GossipLive map[string]int
	Round int
	Sender string
}

type Server struct {
	live map[string]int
	lock sync.Mutex
	Round int
	Address string
	peers []PeerServer
}

func (t *Server) Heartbeat(args *Args,reply *int) error {

}

func (t *Server) sendHeartbeat(to PeerServer) {

}

func (t *Server) GenerateReport() {

}

func main() {

	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	my_address := //TODO: FILL IN
	server.Address = my_address
	server.Round = 0
	server.peers = make([]PeerServer,0)
	server.live = make(map[string]int)
	peer_addresses := []string{/* TODO FILL IN*/}

	time.Sleep(10*time.Second) // WAIT to start other servers

	for _,addr := range(peer_addresses) {
		if(addr == my_address) {
			continue
		}
		client, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		server.peers = append(server.peers, PeerServer{addr,client})
	}

	/*
		TODO: call send heartbeats to a random server every second
			- NOTE: ensure that this code is non-blocking!
		TODO: call generate report every 5 seconds
	*/
	
}
```

Notice that this code sets up an RPC server _and_ an RPC client.

## The Server Struct

```golang
type Server struct {
	live map[string]int
	lock sync.Mutex
	Round int
	Address string
	peers []PeerServer
}
```

The information pertaining to your server is stored in the `Server` struct.

Notice we maintain a logical clock in the variable `Round`.

Notice we are working to maintain a `live` map, which stores peer server's address string mapped to the int round of last contact.

## Sending Heartbeats

Every 1 second you should send out a new heartbeat RPC to one other random server in your group.

- On sending a heartbeat, increment `Round`. 
	- This means, by default, round increases at least once every second.
- The heartbeat message `args` should also forward information about known live servers via sending over a **copy** of our `live` map.
	- Hint: the `maps.Clone()` function that can be imported via the `maps` package may be useful.
	- This forwarding will act as gossip about which servers are known to be alive.

Remember: this is what is included in the arguments to Heartbeat RPC. Populate each field correctly before sending.

```golang
type Args struct {
	GossipLive map[string]int
	Round int
	Sender string
}
```

Ensure that your sending code does not block if one server is not responding!

You will need to include an RPC `Call` that resembles the following:
```
err := to.Client.Call("Server.Heartbeat", args, &reply)
if err != nil {
	log.Println("RPC error:", err)
}
```

> [!IMPORTANT]
> - Complete the implementation of `sendHeartbeat`
> - Add code to `main()` to call `sendHeartbeat` periodically to a random peer **every second**

## Receiving Heartbeats

On receiving a heartbeat:

- If the incoming round > your server's round, update your round!
- Mark the sender as live as of this round in the `live` map
- If the sender has heard from another server in a more recent round, copy this information to your `live` map.


> [!IMPORTANT]
> Complete the implementation of the Heartbeat RPC

## Periodic Reports

Every 5 seconds we want to report which other machines have been contacted in the **past 10 rounds**.


You can format your report how you like, however, a suggestion is below:
```
2026/07/15 09:13:05 REPORT!
2026/07/15 09:13:05 ROUND 117
2026/07/15 09:13:05 map[10.0.0.154:1234:116 10.0.0.154:1235:116 10.0.0.47:1234:117 10.0.0.47:1235:116]
---- LIVE SERVERS ----
10.0.0.47:1234
10.0.0.47:1235
10.0.0.154:1234
10.0.0.154:1235
```

> [!IMPORTANT]
> - Complete the implementation of GenerateReport
> - Add code to `main` to call generateReport **every 5 seconds**.

## Testing

To test your code:
1. populate the `peer_addresses` with the ip addresses of the other students in your group.
2. Run your code with `go run server.go`. Make sure all the other students are also running their servers!

### Finding your IP address
Find the ip address of your current machine.

On Linux:
```
hostname -I
```

On Mac:
```
ipconfig getifaddr en0
```

On Windows:
```
ipconfig
```

Talk to your group mates and write down the IP addresses of their machines. Place these in the servers slice in your `server.go` code. 
**Make sure to include the used port (`:1234`) at the end of each address**

With your group connect all 4-5 laptops together.

> [!IMPORTANT]
> Observe all servers running and reporting live servers.
> Help your classmates debug their code!

### Disconnect a Server

Now, disconnect one laptop from the WiFi-- what happens? How long before they are reported as FAILED?

A bit later, reconnect that laptop to the WiFi-- what happens?

Our desired behavior is that the Live list excludes the wifi disconnected server from the list. Then, when the server is reconnected, it returns to the list.

### Adding Delays

Now, let's add a delay in Heartbeat() RPC before responding.

```
time.Sleep(time.Second * 3)
```

Have one server wait 3 seconds before responding-- is it reported as alive or not?

### Blocking a Specific Computer

One of the students in your group likely has a Mac/Linux computer. They will run an experiment where they block traffic from a machine.

On Mac:
```
sudo pfctl -t blocklist -T add <IP_Address>
```

```
sudo pfctl -t blocklist -T del <IP_Address>
```

Do you reports across machines disagree or agree now that certain computers have blocked traffic?
