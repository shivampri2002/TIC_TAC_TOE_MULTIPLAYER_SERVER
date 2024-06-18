package main

import (
	"log"
)

//in this server app the 0 & 1 will be for player and 2 for tie
//hande move will signal the winner or draw and otherwise 3 for keep going

// HandleMove processes a player's move
func HandleMove(game *Game, player int8, x int, y int) int8 {
	if x < 0 || x >= 3 || y < 0 || y >= 3 {
		log.Println("Invalid Move: index out of bound!")
		return -1
	}
	if game.Board[x][y] != -1 {
		//just check it here while it will not occur as previously checked on recieving msg
		log.Println("cell already occupied")
		return -1
	}

	game.TotalTurn += 1
	game.Board[x][y] = player

	if game.TotalTurn > 3 {
		wn := result(game.Board)
		if wn == 0 || wn == 1 {
			return wn
		}

		//check for tie
		if game.TotalTurn == 9 {
			return 2
		}
	}

	// Switch turn
	game.Turn = switchTurn(player)
	return 3
}

func switchTurn(player int8) int8 {
	if player == 0 {
		return 1
	}
	return 0
}

// func CheckWin(b [3][3]uint8) uint16 {
// 	// Implement win checking logic
// 	winner := result(b)

// 	return 3
// }

// func CheckDraw(board [3][3]uint8) bool {
// 	// Implement draw checking logic
// 	return false
// }

// checking for the any winner in the game
func result(b [3][3]int8) int8 {
	//checking if the top left to bottom right diagonal has any winner
	if b[0][0] != 0 && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return b[0][0]
	}

	//checking if the top right to bottom left diagonal has any winner
	if b[0][2] != 0 && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return b[0][2]
	}

	//checking for the vertical and horizontal line have any winner
	for i := 0; i < 3; i++ {
		if b[i][0] != 0 && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			return b[i][0]
		}

		if b[0][i] != 0 && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
			return b[0][i]
		}
	}

	return -1
}
