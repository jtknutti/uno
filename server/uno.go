package main

import (
	"math/rand"

	"github.com/jak103/uno/db"
	"github.com/jak103/uno/model"
)

//Old Items wont need or use these anymore

// ////////////////////////////////////////////////////////////
// // Utility functions used in place of firebase
// ////////////////////////////////////////////////////////////
// func randColor(i int) string {
// 	switch i {
// 	case 0:
// 		return "red"
// 	case 1:
// 		return "blue"
// 	case 2:
// 		return "green"
// 	case 3:
// 		return "yellow"
// 	}
// 	return ""
// }

// ////////////////////////////////////////////////////////////
// // All the data needed for a simulation of the game
// // eventually, this will be replaced with firebase
// ////////////////////////////////////////////////////////////
// var gameID string = ""
// var currCard []model.Card = nil // The cards are much easier to render as a list
// var players []string = []string{}
// var playerIndex = 0 // Used to iterate through the players
// var currPlayer string = ""
// var allCards map[string][]model.Card = make(map[string][]model.Card) // k: username, v: list of cards
// var gameStarted bool = false

////////////////////////////////////////////////////////////
// Utility functions
////////////////////////////////////////////////////////////
// func newRandomCard() []model.Card {
// TODO use deck utils instead
// 	return []model.Card{model.Card{rand.Intn(10), randColor(rand.Intn(4))}}
// }



func newPayload(user string) map[string]interface{} { // User will default to "" if not passed
	payload := make(map[string]interface{})

	// Return the game model instead of these individually.

	// Update known variables
	// payload["current_card"] = currCard
	// payload["current_player"] = currPlayer
	// payload["all_players"] = players
	// payload["deck"] = allCards[user] // returns nil if currPlayer = "" or user not in allCards
	// payload["game_id"] = gameID
	// payload["game_over"] = checkForWinner()

	return payload
}

func checkID(id string) bool {
	return id == gameID
}

func contains(arr []string, val string) (int, bool) {
	for i, item := range arr {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

////////////////////////////////////////////////////////////
// These are all of the functions for the game -> essentially public functions
////////////////////////////////////////////////////////////
func updateGame(game string, username string) bool {
	success := false
	if success = checkID(game); success && gameStarted {
		return true
	}
	return false
}

func createNewGame() error {
	database, err := db.GetDb()

	if err != nil {
		return err
	}

	game, err := database.CreateGame()

	if err != nil {
		return err
	}

	gameID = game.ID
	return nil
}

func joinGame(game string, username string) bool {
	if checkID(game) {
		user := username

		if _, found := contains(players, user); !found {
			players = append(players, user)
			allCards[user] = nil // No cards yet
		}
		return true
	}
	return false // bad game_id
}

func playCard(game string, username string, card model.Card) bool {
	if checkID(game) && currPlayer == username {
		cards := allCards[username]
		if card.Color == currCard[0].Color || card.Value == currCard[0].Value {
			// Valid card can be played
			playerIndex = (playerIndex + 1) % len(players)
			currPlayer = players[playerIndex]
			currCard[0] = card

			for index, item := range cards {
				if item == currCard[0] {
					allCards[username] = append(cards[:index], cards[index+1:]...)
					break
				}
			}
		}
		return true
	}
	return false
}

// TODO: Keep track of current card that is top of the deck
func drawCard(gameID string, playerID string) bool {
	database, err := db.GetDb()

	if err != nil {
		return err
	}

	game, err := database.LookupGameByID(gameID)
	var player model.Player
	for _, item := range game.Players {
		if playerID == item.ID {
			player = item
		}
	}
	drawCardHelper(game, player)

	database.SaveGame(game)

	return true
}

// TODO: need to deal the actual cards, not just random numbers
func dealCards(game *model.Game) {
	// The game has started, no more players are joining
	// loop through players, set their cards
	gameStarted = true
	currPlayer = players[rand.Intn(len(players))]

	for k := range players {
		cards := []model.Card{}
		for i := 0; i < 7; i++ {

			// TODO Use deck utils instead
			//cards = append(cards, model.Card{rand.Intn(10), randColor(rand.Intn(4))})
			//can set this up similar to the drawCardHelper, but not just call draw card because that makes it call the database too much.
			lastIndex := len(game.DrawPile) - 1
			card := game.DrawPile[lastIndex]
			append(cards, card)
			game.DrawPile = game.DrawPile[:lastIndex]


		}
		allCards[players[k]] = cards
	}

	// TODO Use deck utils instead
	//currCard = newRandomCard()
	//This will draw one more card, but instead of adding it to a players hand it will add it to the discard pile and set it as the current Card
	lastIndex := len(game.DrawPile) - 1
	startCard := game.DrawPile[lastIndex]
	append(DiscardPile.Cards, card)
	game.DrawPile = game.DrawPile[:lastIndex]
	currCard = startCard

}



// TODO: make sure this reflects on the front end
func checkForWinner() string {
	for k := range players {
		if len(allCards[players[k]]) == 0 {
			return players[k]
		}
	}
	return ""
}
