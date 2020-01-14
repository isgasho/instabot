package main

import (
	"context"
	"fmt"
	"log"

	"github.com/winterssy/instabot"
)

func main() {
	bot := instabot.New("YOUR_USERNAME", "YOUR_PASSWORD")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := bot.Login(ctx, true)
	if err != nil {
		log.Fatal(err)
	}

	data, err := bot.GetMe(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}
