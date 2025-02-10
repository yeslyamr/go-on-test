package main

import (
	"goon/internal/bot"
	"goon/internal/storage"
	"goon/internal/stream"
	"log"
)

func main() {
	// DB
	s := storage.New()
	s.Initialize()

	// STREAM HANDLER
	streamHandler := stream.NewHandler(s)
	go streamHandler.Start()
	log.Println("Started stream handler")

	// BOT
	discordBot, err := bot.New(s, streamHandler)
	if err != nil {
		log.Fatalf("Error initializing discord bot: %v", err)
	}

	// START BOT
	err = discordBot.Start()
	if err != nil {
		log.Fatalf("Error starting discord bot: %v", err)
	}

	defer func(discordBot *bot.Bot) {
		_ = discordBot.Close()
	}(discordBot)

	//
	select {}
}
