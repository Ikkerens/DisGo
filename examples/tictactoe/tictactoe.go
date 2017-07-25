package main

import (
	"github.com/ikkerens/disgo"
	"github.com/ikkerens/disgo/examples"
	"github.com/slf4go/logger"
)

var (
	bot   *disgo.User
	games map[disgo.Snowflake]*ticTacToe
)

func init() {
	games = make(map[disgo.Snowflake]*ticTacToe)
}

func main() {
	discord, err := disgo.NewBot(examples.Token)
	if err != nil {
		logger.ErrorE(err)
		return
	}

	discord.RegisterEventHandler(onReady)
	discord.RegisterEventHandler(onMessage)
	discord.RegisterEventHandler(onReactionAdd)

	err = discord.Connect()
	if err != nil {
		logger.ErrorE(err)
		return
	}
	defer discord.Close()

	examples.WaitForExit()
}

func onReady(_ *disgo.Session, event disgo.ReadyEvent) {
	bot = event.User
}

func onMessage(discord *disgo.Session, event disgo.MessageCreateEvent) {
	if len(event.Mentions()) == 2 {
		var (
			valid      bool
			challenged *disgo.User
		)

		for _, mention := range event.Mentions() {
			switch mention.ID() {
			case bot.ID():
				if valid {
					event.Reply("You can't challenge the bot!")
					return
				}
				valid = true
			default:
				challenged = mention
			}
		}

		if valid {
			if event.Author().ID() == challenged.ID() {
				event.Reply("You can't challenge yourself!")
				return
			}
			game := &ticTacToe{discord: discord, player1: event.Author(), player2: challenged}
			if message := game.start(event.ChannelID()); message != 0 {
				games[message] = game
			}
		}
	}
}

func onReactionAdd(discord *disgo.Session, event disgo.MessageReactionAddEvent) {
	if event.UserID != bot.ID() {
		game, exists := games[event.MessageID]
		if exists {
			game.addReaction(event.UserID, event.Emoji.Name())
			discord.MessageDeleteReaction(event.ChannelID, event.MessageID, event.UserID, event.Emoji.Name())
		}
	}
}
