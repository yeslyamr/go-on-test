package bot

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"goon/internal/storage"
	"goon/internal/stream"
	"log"
	"strings"
)

type Bot struct {
	session *discordgo.Session
	storage *storage.Storage
	stream  *stream.Handler
}

func New(storage *storage.Storage, stream *stream.Handler) (*Bot, error) {
	// get the bot token from the environment
	var token string
	flag.StringVar(&token, "token", "", "Discord bot token")
	flag.Parse()
	if token == "" {
		log.Fatal("discord bot token was not provided")
	}

	// create new discord session
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		session: session,
		storage: storage,
		stream:  stream,
	}, nil
}

func (b *Bot) Start() error {
	b.session.AddHandler(b.handler)

	log.Println("Bot is now running")
	return b.session.Open()
}

func (b *Bot) Close() error {
	return b.session.Close()
}

// handler is the main message handler for the bot.
func (b *Bot) handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages that doesn't start with !
	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	// get command and arguments
	args := strings.Fields(m.Content)
	command := args[0]

	// Route the command to the appropriate handler.
	switch command {
	case "!recent":
		b.handleRecentCommand(s, m, args)
	case "!setLang":
		b.handleSetLangCommand(s, m, args)
	case "!stats":
		b.handleStatsCommand(s, m, args)
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Try !recent, !setLang, !stats")
	}
}
