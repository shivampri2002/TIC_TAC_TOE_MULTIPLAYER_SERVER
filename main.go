package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Player represents a player in the game
type Player struct {
	ID   string
	Conn *websocket.Conn
}

// Game represents a Tic Tac Toe game
type Game struct {
	Board     [3][3]int8
	Turn      int8
	TotalTurn uint16
	PlayerX   *Player
	PlayerO   *Player
}

// Games map to hold active games
var Games = make(map[string]*Game)
var pendingPlayer *Player
var mu sync.Mutex

func main() {
	app := fiber.New()

	// WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		player := &Player{
			ID:   uuid.New().String(),
			Conn: c,
		}
		log.Println("Joined Player: ", player)

		mu.Lock()
		if pendingPlayer == nil {
			// Wait for another player to connect
			pendingPlayer = player
			mu.Unlock()
			c.WriteMessage(websocket.TextMessage, []byte("WAIT"))
		} else {
			// Start a new game
			gameID := uuid.New().String()
			game := &Game{
				Turn:      0,
				TotalTurn: 0,
				PlayerX:   pendingPlayer,
				PlayerO:   player,
			}
			Games[gameID] = game

			//set the board to -1
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					game.Board[i][j] = -1
				}
			}

			// Reset pending player
			pendingPlayer = nil
			mu.Unlock()

			// Notify players
			// pendingPlayer.Conn.WriteMessage(websocket.TextMessage, []byte("Game started! You are X. Game ID: "+gameID))
			// player.Conn.WriteMessage(websocket.TextMessage, []byte("Game started! You are O. Game ID: "+gameID))

			handleGame(game)

			if pendingPlayer.Conn != nil {
				pendingPlayer.Conn.Close()
			}
			if player.Conn != nil {
				player.Conn.Close()
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}

func handleGame(game *Game) {
	game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("MOVE 10 10"))

	for {
		if game.Turn == 0 {
			_, msg1, err := game.PlayerO.Conn.ReadMessage()
			if err != nil {
				log.Println("read1: ", err)
				break
			}

			mov1 := strings.Split(strings.TrimSpace(string(msg1)), " ")
			if len(mov1) < 3 || mov1[0] != "MOVED" {
				log.Println("Less no of arguments or wrong arguments recieved in moves", mov1)
				return
			}

			x1, err := strconv.Atoi(mov1[1])
			if err != nil {
				log.Println("Wrong moved index: ", mov1)
				return
			}
			y1, err := strconv.Atoi(mov1[2])
			if err != nil {
				log.Println("Wrong moved index: ", mov1)
				return
			}

			res1 := HandleMove(game, game.Turn, x1, y1)

			if res1 == -1 {
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("Unauthorized Move"))
				log.Println("Unauthorized move made by the playerO")
				return
			} else if res1 == 0 {
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER W"))
				game.PlayerX.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER L"))
				return
			} else if res1 == 1 {
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER L"))
				game.PlayerX.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER W"))
				return
			} else if res1 == 2 {
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER T"))
				game.PlayerX.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER T"))
				return
			}

			game.PlayerX.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("MOVE %d %d", x1, y1)))
		} else {
			_, msg2, err := game.PlayerX.Conn.ReadMessage()
			if err != nil {
				log.Println("read2: ", err)
			}

			mov2 := strings.Split(strings.TrimSpace(string(msg2)), " ")
			if len(mov2) < 3 || mov2[0] != "MOVED" {
				log.Println("Less no of arguments or wrong arguments recieved in moves", mov2)
				return
			}

			x2, err := strconv.Atoi(mov2[1])
			if err != nil {
				log.Println("Wrong moved index: ", mov2)
				return
			}
			y2, err := strconv.Atoi(mov2[2])
			if err != nil {
				log.Println("Wrong moved index: ", mov2)
				return
			}

			res2 := HandleMove(game, game.Turn, x2, y2)

			if res2 == -1 {
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("Unauthorized Move"))
				log.Println("Unauthorized move made by the playerO")
				return
			} else if res2 == 0 {
				game.PlayerX.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER L"))
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER W"))
				return
			} else if res2 == 1 {
				game.PlayerX.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER W"))
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER L"))
				return
			} else if res2 == 2 {
				game.PlayerX.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER T"))
				game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte("GAME_OVER T"))
				return
			}

			game.PlayerO.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("MOVE %d %d", x2, y2)))
		}

		//on getting the response from the user check if they are making valid move if not send error
		// Ensure the player is authorized to make a move in this game
		// if (game.PlayerX.ID == playerID && game.Turn == "X") || (game.PlayerO.ID == playerID && game.Turn == "O") {
		// 	// Process the move (you need to implement game logic)
		// 	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		// 		log.Println("write:", err)
		// 		break
		// 	}
		// } else {
		// 	conn.WriteMessage(websocket.TextMessage, []byte("Not your turn or unauthorized"))
		// }
	}
}
