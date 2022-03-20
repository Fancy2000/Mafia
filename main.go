package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"hw5/service"
	"log"
	"math/rand"
	"net"
	"time"
)

type Users struct {
	Role   string
	Name   string
	Status string
}

func (Users) Death(this Users) Users {
	this.Status = "dead"
	return this
}

//var Players []string
var Players = make(map[string]Users)
var Votion = make(map[string]int)
var CountPlayersEndDay int32 = 0
var CountPlayersEndNight int32 = 0
var CountPlayersVote int32 = 0
var CommissionerFindMafia = false
var left int32 = 4

type server struct {
	service.UnimplementedCommunicationServer
}

var Roles = []string{"Mafia", "commissioner", "civilian", "civilian"}

func (server) ListOfPlayers(ctx context.Context, empty *service.Empty) (*service.ResponseListOfPlayers, error) {
	map_info := make(map[string]string)
	for key, value := range Players {
		map_info[key] = "Player: " + value.Name + " is " + value.Status
		if value.Status == "dead" {
			map_info[key] += " He was a " + value.Role
		}
	}
	res := service.ResponseListOfPlayers{
		Members: map_info,
	}
	return &res, nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	if min > max {
		return min
	} else {
		return rand.Intn(max-min) + min
	}
}

func (server) CountAlivePeople(ctx context.Context, empty *service.Empty) (*service.HowMany, error) {
	res := service.HowMany{People: left}
	return &res, nil
}

func (server) EndGame(ctx context.Context, empty *service.Empty) (*service.Request, error) {
	mafia_alive := true
	fmt.Println(Players)
	for _, val := range Players {
		if val.Role == "Mafia" {
			if val.Status == "dead" {
				mafia_alive = false
			}
		}
	}
	if !mafia_alive {
		res := service.Request{Name: "Mafia is dead. the civilians WON!!!"}
		fmt.Println("Mafia is dead. the civilians WON!!!")
		return &res, nil
	}
	fmt.Println("left: ", left)
	if left <= 2 {
		res := service.Request{Name: "Mafia WON!!!"}
		fmt.Println("Mafia WON!!!")

		return &res, nil
	}
	res := service.Request{Name: "continue"}
	return &res, nil
}

func (server) KillAfterVote(ctx context.Context, empty *service.Empty) (*service.Request, error) {
	max := 0
	for _, val := range Votion {
		if val > max {
			max = val
		}
	}
	name := ""
	for key, val := range Votion {
		if max == val {
			name = key
		}
	}
	for key, _ := range Votion {
		Votion[key] = 0
	}
	Players[name] = Players[name].Death(Players[name])
	left -= 1
	res := service.Request{Name: name + "is killed. He was a " + Players[name].Role}
	return &res, nil
}

func (server) CountVote(ctx context.Context, req *service.Empty) (*service.HowMany, error) {
	res := service.HowMany{People: CountPlayersVote}
	return &res, nil
}

func (server) Vote(ctx context.Context, req *service.Votion) (*service.Request, error) {
	_, find := Players[req.Name]
	if !find || Players[req.Name].Status == "dead" {
		res := service.Request{Name: "Error"}
		return &res, nil
	}
	CountPlayersVote += 1
	fmt.Println(CountPlayersVote)
	if CountPlayersVote == req.People {
		CountPlayersVote = 0
	}
	res := service.Request{Name: "OK"}
	Votion[req.Name] += 1
	return &res, nil
}

func (server) CheckDeadAfterNight(ctx context.Context, empty *service.Empty) (*service.Request, error) {
	ans := ""
	for key, val := range Players {
		if val.Status == "dead" {
			ans += key + " "
		}
	}
	res := service.Request{Name: ans}
	return &res, nil
}

func (server) CheckWhoMafia(ctx context.Context, empty *service.Empty) (*service.Request, error) {
	if CommissionerFindMafia {
		for _, val := range Players {
			if val.Role == "Mafia" {
				res := service.Request{Name: "Mafia is: " + val.Name}
				return &res, nil
			}
		}
	}
	res := service.Request{Name: "NO"}
	return &res, nil
}

func (server) MafiaKill(ctx context.Context, req *service.Request) (*service.Request, error) {
	_, find := Players[req.Name]
	if !find || Players[req.Name].Status == "dead" {
		res := service.Request{Name: "Error"}
		return &res, nil
	}
	res := service.Request{
		Name: "OK",
	}
	Players[req.Name] = Players[req.Name].Death(Players[req.Name])
	fmt.Println("Mafia kill: ", req.Name)
	left -= 1
	return &res, nil
}

func (server) CommissionerCheck(ctx context.Context, req *service.Request) (*service.Request, error) {
	_, find := Players[req.Name]
	if !find || Players[req.Name].Status == "dead" {
		res := service.Request{Name: "Error"}
		return &res, nil
	}
	res := service.Request{Name: Players[req.Name].Role}
	if Players[req.Name].Role == "Mafia" {
		CommissionerFindMafia = true
	}
	return &res, nil
}

func (server) Update(ctx context.Context, req *service.Request) (*service.Request, error) {
	if Players[req.Name].Status == "dead" {
		res := service.Request{Name: "dead"}
		return &res, nil
	} else {
		res := service.Request{Name: "alive"}
		return &res, nil
	}
}

func (server) DeadVote(ctx context.Context, empty *service.Empty) (*service.Empty, error) {
	CountPlayersVote += 1
	fmt.Println(CountPlayersVote)
	res := service.Empty{}
	return &res, nil
}

func (server) SayHwoIsMafia(ctx context.Context, req *service.Empty) (*service.Empty, error) {
	for key, val := range Players {
		if val.Role == "Mafia" {
			fmt.Println("MAFIA IS: " + key)
		}
	}
	res := service.Empty{}
	return &res, nil
}

func (server) CheckHowManyPeopleIsNotSleep(ctx context.Context, empty *service.Empty) (*service.HowMany, error) {
	res := service.HowMany{People: CountPlayersEndDay}
	return &res, nil
}

func (server) CheckHowManyPeopleIsSleep(ctx context.Context, empty *service.Empty) (*service.HowMany, error) {
	res := service.HowMany{People: CountPlayersEndNight}
	return &res, nil
}

func (server) StartGame(ctx context.Context, empty *service.Empty) (*service.HowMany, error) {
	res := service.HowMany{People: int32(len(Players))}
	return &res, nil
}

func (server) EndNight(ctx context.Context, req *service.HowMany) (*service.HowMany, error) {
	CountPlayersEndNight += 1
	res := service.HowMany{People: CountPlayersEndNight}
	if CountPlayersEndNight == req.People {
		CountPlayersEndNight = 0
	}
	return &res, nil
}

func (server) EndDay(ctx context.Context, empty *service.Empty) (*service.HowMany, error) {
	CountPlayersEndDay += 1
	res := service.HowMany{People: CountPlayersEndDay}
	if CountPlayersEndDay == left {
		CountPlayersEndDay = 0
	}
	return &res, nil
}

func (server) InitPlayer(ctx context.Context, req *service.Request) (*service.ResponseInit, error) {
	if len(Players) == 4 {

		res := service.ResponseInit{
			Name: "all places are occupied",
		}
		return &res, errors.New("all places are occupied")
	}
	_, found := Players[req.Name]
	dublicated_name := false
	if found {
		dublicated_name = true
		req.Name = req.Name + "1"
		for {
			_, found = Players[req.Name]
			if !found {
				break
			}
			req.Name = req.Name + "1"
		}
	}
	var index int
	if len(Roles) == 1 {
		index = 0
	} else {
		index = random(0, len(Roles)-1)
	}

	Players[req.Name] = Users{Roles[index], req.Name, "alive"}
	Roles[index] = Roles[len(Roles)-1]
	Roles[len(Roles)-1] = ""
	Roles = Roles[:len(Roles)-1]

	var ans string
	if dublicated_name {
		ans = "Sorry your name was taken so we added 1 to it.\n" +
			"Ok " + req.Name + " your role is " + Players[req.Name].Role
	} else {
		ans = "Ok " + req.Name + " your role is " + Players[req.Name].Role
	}
	res := service.ResponseInit{
		Msg:  ans,
		Name: req.Name,
		Role: Players[req.Name].Role,
	}
	fmt.Println("Player " + res.Name + " connected")
	return &res, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:5050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service.RegisterCommunicationServer(s, &server{})
	fmt.Println("Server start")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
