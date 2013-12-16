package main

import (
	"log"
	"flag"
	"strings"
	"time"

	"github.com/hashicorp/serf/command/agent"

	"github.com/vito/maestro/warden_pool"
)

var serfAgentRPCAddr = flag.String(
	"serfAgentRPCAddr",
	"127.0.0.1:7373",
	"local serf agent's RPC address",
)

var serfMembers = flag.String(
	"serfMembers",
	"",
	"join a serf cluster",
)

func main() {
	flag.Parse()

	var serfClient *agent.RPCClient
	var err error

	interval := 1 * time.Second

	for {
		serfClient, err = agent.NewRPCClient(*serfAgentRPCAddr)
		if err == nil {
			break
		}

		log.Println(
			"failed to reach serf agent at",
			*serfAgentRPCAddr,
			"trying again in",
			interval,
		)

		time.Sleep(interval)
	}

	pool := warden_pool.New(serfClient)

	err = pool.Listen()
	if err != nil {
		log.Fatalln("failed to listen for events:", err)
	}

	addrs := strings.Split(*serfMembers, ",")
	joinedCount, err := serfClient.Join(addrs, false)
	if err != nil {
		log.Fatalln("failed to join serf cluster:", err)
	}

	log.Println("joined", joinedCount, "members of serf cluster")

	select {}
}
