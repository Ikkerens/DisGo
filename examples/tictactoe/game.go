package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/ikkerens/disgo"
	"github.com/slf4go/logger"
)

const (
	expireTime = 60

	slotKeyCap = '\U000020e3'
	slotCircle = '\U0001f534'
	slotCross  = '\U0000274e'
)

var winConditions = [...][]int{
	// Horizontal
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},

	// Vertical
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},

	// Diagonal
	{0, 4, 8},
	{2, 4, 6},
}

type ticTacToe struct {
	discord                 *disgo.Session
	player1, player2        *disgo.User
	gameOver, draw, expired bool

	board *disgo.Message
	slots [9]rune
	turn  bool

	expire *time.Timer
}

func (game *ticTacToe) start(channel disgo.Snowflake) disgo.Snowflake {
	for i := range game.slots {
		game.slots[i] = 0
	}

	board, err := game.discord.SendEmbed(channel, game.buildBoard())
	if err != nil {
		logger.ErrorE(err)
		return 0
	}
	game.board = board
	game.expire = time.NewTimer(expireTime * time.Second)

	go func() {
		for i := 1; i <= 9; i++ {
			board.AddReaction(strconv.Itoa(i) + string(slotKeyCap))
		}
	}()

	go func() {
		<-game.expire.C
		game.expired = true
		game.end()
	}()

	return board.ID()
}

func (game *ticTacToe) addReaction(user disgo.Snowflake, emoji string) {
	current, _, piece := game.getCurrentPlayer()
	if current.ID() != user {
		return
	}

	runes := []rune(emoji)
	if utf8.RuneCountInString(emoji) == 2 && runes[1] == slotKeyCap {
		slot, err := strconv.ParseInt(string(runes[0]), 10, 8)
		if err != nil {
			return
		}

		if slot >= 1 && slot <= 9 && game.slots[slot-1] == 0 {
			game.expire.Reset(expireTime * time.Second)
			game.slots[slot-1] = piece
			game.turn = !game.turn

			if !game.checkForEndOfGame() {
				game.board.EditEmbed(*game.buildBoard())
				game.board.DeleteOwnReaction(emoji)
			}
		}
	}
}

func (game *ticTacToe) checkForEndOfGame() bool {
	// Is there a winner?
	for _, win := range winConditions {
		if game.slots[win[0]] == game.slots[win[1]] && game.slots[win[0]] == game.slots[win[2]] && game.slots[win[0]] != 0 {
			if game.slots[win[0]] == slotCircle {
				game.turn = false
			} else {
				game.turn = true
			}
			game.end()
			return true
		}
	}

	// Are there free slots?
	for _, slot := range game.slots {
		if slot == 0 {
			return false
		}
	}

	// No winner, but also no free slots, must be a draw
	game.draw = true
	game.end()
	return true
}

func (game *ticTacToe) end() {
	if !game.gameOver {
		delete(games, game.board.ID())
		game.expire.Stop()
		game.gameOver = true
		game.board.DeleteAllReactions()
		game.board.EditEmbed(*game.buildBoard())
	}
}

func (game *ticTacToe) buildBoard() *disgo.Embed {
	var field bytes.Buffer

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			slot := game.slots[y*3+x]
			if slot == 0 {
				field.WriteString(strconv.Itoa(y*3 + x + 1))
				field.WriteRune(slotKeyCap)
			} else {
				field.WriteRune(slot)
			}
		}

		if y != 2 {
			field.WriteRune('\n')
		}
	}

	var footer string
	current, color, _ := game.getCurrentPlayer()

	if game.gameOver {
		color = 0x4F545C
		if game.draw {
			footer = "The game ended in a draw."
		} else if game.expired {
			footer = "This game has been cancelled due to player inactivity."
		} else {
			footer = fmt.Sprintf("%s won the game.", current.Username())
		}
	} else {
		footer = fmt.Sprintf("It's currently %s's turn.", current.Username())
	}

	return &disgo.Embed{
		Title:       fmt.Sprintf("Tic Tac Toe: %c %s vs %c %s", slotCircle, game.player1.Username(), slotCross, game.player2.Username()),
		Description: field.String(),
		Color:       color,
		Footer: disgo.EmbedFooter{
			Text:    footer,
			IconURL: current.AvatarURL(),
		},
	}
}

func (game *ticTacToe) getCurrentPlayer() (*disgo.User, int, rune) {
	if game.turn {
		return game.player2, 0x77B255, slotCross
	} else {
		return game.player1, 0xDD2E44, slotCircle
	}
}
