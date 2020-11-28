package player

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Player ...
type Player struct {
	Id    int
	Conn  clientFrontEnd
	Name  string
	Msgs  chan string
	PList PlayerList
	comm  communicator
}

func (p *Player) Close(msg string) {
	log.Printf("%s\n", msg)
	delete(p.PList, p.Id)
	p.Conn.Close()
}

func (p Player) GetId() int {
	return p.Id
}

func (p Player) GetName() string {
	return p.Name
}

func SetupNewPlayer(conn net.Conn, id int, PlayerList PlayerList, comm communicator) {
	var p *Player
	var msgs = make(chan string)
	p = &Player{
		Id:    id,
		Conn:  conn,
		Msgs:  msgs,
		PList: PlayerList,
		comm:  comm,
	}
	err := getPlayerName(p)
	if err != nil {
		p.Close(err.Error())
		return
	}
	PlayerList[p.Id] = p
	p.comm.SendMsgTo(fmt.Sprintf("You are Player %d", id), *p)
	go p.comm.ListenForMessages(*p)
	go p.comm.EchoMessages(*p, PlayerList)
}

func getPlayerName(p *Player) error {
	p.comm.SendMsgTo("Hello! What is your name? ", *p)
	reader := bufio.NewReader(p.Conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	p.Name = name[:len(name)-1]
	return nil
}
