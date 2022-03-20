package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"hw5/service"
	"log"
	"time"
)

var commissioner_check_mafia string
var FindMafia = false
var AlivePeople int32 = 0

type Player struct {
	Role   string
	Name   string
	Status string
}

func info(client service.CommunicationClient) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	res, err := client.ListOfPlayers(ctx, &service.Empty{})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("------------------------")
	for _, value := range res.Members {
		fmt.Println(value)
	}
	fmt.Println("------------------------")
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:5050", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	var name string
	fmt.Println("Write your name:")
	fmt.Scanf("%s\n", &name)
	client := service.NewCommunicationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("Client start")

	//new player
	resp, err := client.InitPlayer(ctx, &service.Request{Name: name})
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(resp.Msg)
	player := Player{Role: resp.Role, Name: resp.Name}
	//wait other players
	fmt.Println("Wait for other players")
	for {

		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		res, err := client.StartGame(ctx, &service.Empty{})
		if err != nil {
			log.Println(err)
		}
		if res.People == 4 {
			break
		}
		time.Sleep(4 * time.Second)
	}
	fmt.Println("Game is Start!")
	fmt.Println("it's a day, get to know each other when you're done enter \"night\"")
	for {
		fmt.Println("Write a command")
		var command string
		fmt.Scanf("%s\n", &command)
		if command == "night" {
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err := client.EndDay(ctx, &service.Empty{})
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Wait for other players")
			resp, err := client.CountAlivePeople(ctx, &service.Empty{})
			AlivePeople = resp.People
			for {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				res, err := client.CheckHowManyPeopleIsNotSleep(ctx, &service.Empty{})
				if err != nil {
					log.Println(err)
				}
				if res.People == AlivePeople || res.People == 0 {
					break
				}
				time.Sleep(4 * time.Second)
			}
			if player.Role == "commissioner" && player.Status != "dead" {
				fmt.Println("write the name of the person you want to check. You can check only alive people. Possible names:")
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				res, err := client.ListOfPlayers(ctx, &service.Empty{})
				if err != nil {
					log.Println(err)
				}
				fmt.Println("------------------------")
				for _, value := range res.Members {
					fmt.Println(value)
				}
				fmt.Println("------------------------")
				var check_name string
				fmt.Scanf("%s\n", &check_name)
				if check_name == player.Name {
					fmt.Println("you can't write your name. Try again")
					for {
						fmt.Scanf("%s\n", &check_name)
						if check_name != player.Name {
							break
						}
						fmt.Println("you can't write your name. Try again")
					}
				}
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				ans, err := client.CommissionerCheck(ctx, &service.Request{Name: check_name})
				if err != nil {
					log.Println(err)
				}
				if ans.Name == "Error" {
					fmt.Println("You can't check this player. Please try again")
					for {
						ctx, cancel = context.WithTimeout(context.Background(), time.Second)
						var check_name_in_error string
						fmt.Scanf("%s\n", &check_name_in_error)
						ans, err := client.CommissionerCheck(ctx, &service.Request{Name: check_name_in_error})
						if err != nil {
							log.Println(err)
						}
						if ans.Name != "Error" {
							if ans.Name == "Mafia" {
								commissioner_check_mafia = ans.Name
								FindMafia = true
							}
							break
						}
						fmt.Println("You can't check this player. Please try again")
					}
				}

				fmt.Println("This player is: " + ans.Name)
				if ans.Name == "Mafia" {
					commissioner_check_mafia = ans.Name
				}

			} else if player.Role == "Mafia" && player.Status != "dead" {
				fmt.Println("write the name of the person you want to Kill. You can Kill only alive people. Possible names:")
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				res, err := client.ListOfPlayers(ctx, &service.Empty{})
				if err != nil {
					log.Println(err)
				}
				fmt.Println("------------------------")
				for _, value := range res.Members {
					fmt.Println(value)
				}
				fmt.Println("------------------------")
				var check_name string
				fmt.Scanf("%s\n", &check_name)
				if check_name == player.Name {
					fmt.Println("you can't write your name. Try again")
					for {
						fmt.Scanf("%s\n", &check_name)
						if check_name != player.Name {
							break
						}
						fmt.Println("you can't write your name. Try again")
					}
				}
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				ans, err := client.MafiaKill(ctx, &service.Request{Name: check_name})
				if err != nil {
					log.Println(err)
				}
				if ans.Name == "Error" {
					fmt.Println("You can't Kill this player. Please try again")
					for {
						ctx, cancel = context.WithTimeout(context.Background(), time.Second)
						var check_name_in_error string
						fmt.Scanf("%s\n", &check_name_in_error)
						ans, err := client.MafiaKill(ctx, &service.Request{Name: check_name_in_error})
						if err != nil {
							log.Println(err)
						}
						if ans.Name != "Error" {
							fmt.Println("You Kill: " + check_name_in_error)
							break
						}
						fmt.Println("You can't Kill this player. Please try again")
					}
				}
				fmt.Println("You Kill: " + check_name)

			}
			fmt.Println("Write \"Day\" to wake up")
			fmt.Scanf("%s\n", &command)
			if command != "Day" {
				fmt.Println("Incorrect command. Please try again")
				for {
					fmt.Scanf("%s\n", &command)
					if command == "Day" {
						break
					}
					fmt.Println("Incorrect command. Please try again")
				}
			}
			//wake up
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.EndNight(ctx, &service.HowMany{People: AlivePeople})
			if err != nil {
				log.Println(err)
			}

			fmt.Println("Wait for other players")
			for {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				sleep, err := client.CheckHowManyPeopleIsSleep(ctx, &service.Empty{})
				if err != nil {
					log.Println(err)
				}
				if sleep.People == AlivePeople || sleep.People == 0 {
					break
				}
				time.Sleep(4 * time.Second)
			}

			// DAY

			//update status
			ans, err := client.Update(ctx, &service.Request{Name: player.Name})
			if err != nil {
				log.Println(err)
			}
			if ans.Name == "dead" {
				player.Status = "dead"
			}

			fmt.Println("OKEY this is a results after night:")
			info(client)

			if player.Role == "commissioner" {
				fmt.Println("You can say hwo is Mafia. To accept write \"accept\"")
				var command string
				fmt.Scanf("%s", &command)
				if command == "accept" {
					if FindMafia {
						ctx, cancel = context.WithTimeout(context.Background(), time.Second)
						_, err := client.CheckDeadAfterNight(ctx, &service.Empty{})
						if err != nil {
							log.Println(err)
						}
					} else {
						fmt.Println("You didn't find Mafia yet")
					}
				}
			}
			//check hwo mafia
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			ans, err = client.CheckDeadAfterNight(ctx, &service.Empty{})
			if err != nil {
				log.Println(err)
			}
			if ans.Name != "NO" {
				fmt.Println(ans)
			}
			//vote
			if player.Status != "dead" {
				fmt.Println("Now you have to vote who is a Mafia")
				fmt.Println("Write your assumption")
				vote_name := ""
				fmt.Scanf("%s\n", &vote_name)
				if vote_name == player.Name {
					fmt.Println("You can't vote for yourself. Please try again")
					for {
						fmt.Scanf("%s\n", &vote_name)
						if vote_name != player.Name {
							break
						}
						fmt.Println("You can't vote for yourself. Please try again")
					}
				}
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				ans, err = client.Vote(ctx, &service.Votion{Name: vote_name, People: AlivePeople})
				if err != nil {
					log.Println(err)
				}
				if ans.Name == "Error" {
					fmt.Println("Incorrect Vote. Please try again")
					for {
						ctx, cancel = context.WithTimeout(context.Background(), time.Second)
						var vote_name_in_error string
						fmt.Scanf("%s\n", &vote_name_in_error)
						ctx, cancel = context.WithTimeout(context.Background(), time.Second)
						ans, err = client.Vote(ctx, &service.Votion{Name: vote_name_in_error, People: AlivePeople})
						if err != nil {
							log.Println(err)
						}
						if ans.Name != "Error" {
							break
						}
						fmt.Println("Incorrect Vote. Please try again")
					}
				}
			} else {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				_, err := client.DeadVote(ctx, &service.Empty{})
				if err != nil {
					log.Println(err)
				}
			}
			//wait other vote
			fmt.Println("Wait for other players")
			for {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				vote, err := client.CountVote(ctx, &service.Empty{})
				if err != nil {
					log.Println(err)
				}
				if vote.People == AlivePeople || vote.People == 0 {
					break
				}
				time.Sleep(4 * time.Second)
			}

			//kill after vote
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			ans, err = client.KillAfterVote(ctx, &service.Empty{})
			if err != nil {
				log.Println(err)
			}
			fmt.Println(ans.Name)

			//end game
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			ans, err = client.EndGame(ctx, &service.Empty{})
			if err != nil {
				log.Println(err)
			}
			if ans.Name != "continue" {
				fmt.Println(ans.Name)
				break
			}
			fmt.Println("When you're ready to sleep enter \"night\"")

		} else if command == "info" {
			info(client)
		} else {
			fmt.Println("Unknown command. Please try again")
		}

	}
}
