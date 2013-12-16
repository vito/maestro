package warden_pool

import (
	"fmt"
	"encoding/json"
	"sync"
	"log"

	"github.com/hashicorp/serf/command/agent"
)

type WardenPool struct {
	serfClient *agent.RPCClient

	members map[string]WardenMember

	eventStream chan map[string]interface{}
	listenHandle agent.StreamHandle

	sync.RWMutex
}

type WardenMember struct {
	Addr string
	AvailableMemory uint64
	AvailableDisk uint64
}

func New(serfClient *agent.RPCClient) *WardenPool {
	return &WardenPool{
		serfClient: serfClient,

		members: make(map[string]WardenMember),

		eventStream: make(chan map[string]interface{}),
	}
}

func (p *WardenPool) Listen() error {
	handle, err := p.serfClient.Stream("user,member-leave", p.eventStream)
	if err != nil {
		return err
	}

	go p.handleEvents()

	p.listenHandle = handle

	return nil
}

func (p *WardenPool) handleEvents() {
	for {
		event, ok := <-p.eventStream
		if !ok {
			break
		}

		log.Println("got event:", event)

		switch event["Event"] {
		case "member-leave":
			members, ok := event["Members"].([]map[string]interface{})
			if !ok {
				continue
			}

			for _, member := range members {
				if member["Role"] != "warden" {
					continue
				}

				addr, ok := member["Addr"].([]interface{})
				if !ok {
					continue
				}

				name := fmt.Sprintf("%d.%d.%d.%d", addr...)

				p.removeMember(name)
			}

		case "user":
			payload := event["Payload"].([]byte)

			switch event["Name"] {
			case "warden.capacity":
				var member WardenMember

				err := json.Unmarshal(payload, &member)
				if err != nil {
					log.Println("invalid warden.capacity:", err)
					continue
				}

				p.addMember(member)
			}
		}
	}
}

func (p *WardenPool) removeMember(addr string) {
	log.Println("removing member:", addr)

	p.Lock()
	defer p.Unlock()

	delete(p.members, addr)
}

func (p *WardenPool) addMember(member WardenMember) {
	log.Println("adding member:", member)

	p.Lock()
	defer p.Unlock()

	p.members[member.Addr] = member
}
