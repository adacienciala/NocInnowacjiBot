package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

const (
	BOT_TOKEN   = "YOUR_TOKEN"
	PREFIX      = "go"
	KITCHEN_URL = "https://pizza-ni.herokuapp.com/api/kitchen"
)

type Pizza struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type PizzaResp struct {
	Name   string  `json:"name"`
	Amount int     `json:"amount"`
	Price  float32 `json:"price"`
	Error  string  `json:"error,omitempty"`
}

var orders []Pizza

func takeOrder(pizza Pizza) float32 {
	reqBody, _ := json.Marshal(pizza)
	reqBytes := bytes.NewBuffer(reqBody)
	resp, err := http.Post(KITCHEN_URL, "application/json", reqBytes)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var respPizza PizzaResp
	err = json.Unmarshal(body, &respPizza)
	if err != nil {
		log.Fatalf("An Error Occured %v", resp.Body)
	}

	if respPizza.Error != "" {
		log.Fatalf("Error w kuchni %v", respPizza.Error)
	}

	return respPizza.Price
}

func handleMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID {
		return
	}
	fmt.Println("Got a message, ", msg.Content)
	message := strings.Split(strings.ToLower(msg.Content), " ")

	if message[0] == PREFIX {
		commands := message[1:]
		fmt.Println("Got command, ", commands)

		// YOUR CODE
		switch commands[0] {
		case "help":
			session.ChannelMessageSend(msg.ChannelID, "I'll look for therapy places for you in my free time")
		case "make":
			orderCount, _ := strconv.Atoi(commands[1])
			for i := 1; i <= orderCount; i = i + 1 {
				pizza := Pizza{Name: commands[2], Amount: orderCount}
				orders = append(orders, pizza)
				message := fmt.Sprintf("%s pizza no %d in the making", pizza.Name, i)
				session.ChannelMessageSend(msg.ChannelID, message)
			}
		case "order":
			orderCount, _ := strconv.Atoi(commands[1])
			pizza := Pizza{Name: commands[2], Amount: orderCount}
			message := fmt.Sprintf("%d x %s pizza in the making", pizza.Amount, pizza.Name)
			session.ChannelMessageSend(msg.ChannelID, message)
			price := takeOrder(pizza)
			message = fmt.Sprintf("Done! It'll be %f", price)
			session.ChannelMessageSend(msg.ChannelID, message)
		case "orders":
			message := fmt.Sprintf("All orders: %v", orders)
			session.ChannelMessageSend(msg.ChannelID, message)
		case "clear":
			orders = nil
			session.ChannelMessageSend(msg.ChannelID, "Orders cleared")
		}
		//
	}
}

func createHttpServer() {
	router := httprouter.New()
	port, present := os.LookupEnv("PORT")
	if !present {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	bot, err := discordgo.New("Bot " + BOT_TOKEN)
	fmt.Println("API version:", discordgo.APIVersion)
	if err != nil {
		fmt.Println("Error creating bot session!")
		panic(err)
	}
	bot.AddHandler(handleMessage)
	bot.Open()
	createHttpServer()
	bot.Close()
}
