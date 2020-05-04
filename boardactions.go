package main

import (
	"fmt"
	"math"
	"math/rand"
)

type node struct {
	wins int
	sims int

	boardState [][]square

	children  []*node
	move      piecemove
	parent    *node
	movecolor string
}

type move struct {
	X             int         `json:"x"`
	Y             int         `json:"y"`
	Flag          string      `json:"flag"`
	PositionOfOpp opponentPos `json:"positionOfOpp"`
}

type piecemove struct {
	Piece     piece `json:"piece"`
	Piecemove move  `json:"pieceMove"`
}

type piece struct {
	Xpos int `json:"xpos"`
	Ypos int `json:"ypos"`
}

type opponentPos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func montecarlomove(board [][]square, color string) piecemove {

	var root node
	root.boardState = board
	root.sims = 0
	root.wins = 0
	root.children = nil
	root.movecolor = opposingColor(color)

	var startingMoves = generateAllValidMoves(board, color)

	if len(startingMoves) == 1 {
		return startingMoves[0]
	}

	for _, move := range startingMoves {

		root.children = append(root.children, createNewNode(move, board, color, &root))

	}

	for i := 0; i < 1000; i++ {

		var promisingNode = selectNode(&root)
		expandNode(promisingNode)

		var testNode = promisingNode

		if len(promisingNode.children) > 0 {
			testNode = promisingNode.children[(int)(math.Floor((float64)(rand.Intn(len(promisingNode.children)))))]

		}

		playOut(testNode, color)

	}

	return selectBestMove(root)
}

func selectBestMove(root node) piecemove {
	var bestNode = root.children[0]

	for _, node := range root.children {
		if bestNode.sims < node.sims {
			bestNode = node
		}
	}

	return bestNode.move
}

func selectNode(node *node) *node {

	if len(node.children) > 0 {
		mostValue := node.children[0]

		for _, child := range node.children {
			if uctValue(*child) > uctValue(*mostValue) {
				mostValue = child
			}

		}

		return selectNode(mostValue)

	}

	return node

}

func createNewNode(movetomake piecemove, board [][]square, movecolor string, parent *node) *node {
	var newNode node
	newNode.boardState = playmove(movetomake, copyBoard(board))
	newNode.movecolor = movecolor
	newNode.sims = 0
	newNode.wins = 0
	newNode.move = movetomake
	newNode.parent = parent
	return &newNode
}

func playmove(movetomake piecemove, board [][]square) [][]square {

	var newBoard = board

	var pColor = newBoard[movetomake.Piece.Xpos][movetomake.Piece.Ypos].Color
	var pKing = newBoard[movetomake.Piece.Xpos][movetomake.Piece.Ypos].King

	newBoard[movetomake.Piece.Xpos][movetomake.Piece.Ypos].Color = "null"
	newBoard[movetomake.Piece.Xpos][movetomake.Piece.Ypos].King = false

	newBoard[movetomake.Piecemove.X][movetomake.Piecemove.Y].Color = pColor
	newBoard[movetomake.Piecemove.X][movetomake.Piecemove.Y].King = pKing

	if movetomake.Piecemove.Flag == "capture" {
		newBoard[movetomake.Piecemove.PositionOfOpp.X][movetomake.Piecemove.PositionOfOpp.Y].Color = "null"
		newBoard[movetomake.Piecemove.PositionOfOpp.X][movetomake.Piecemove.PositionOfOpp.Y].King = false
	}

	if (movetomake.Piecemove.X == 0 && pColor == "w") ||
		(movetomake.Piecemove.X == 7 && pColor == "r") {
		newBoard[movetomake.Piecemove.X][movetomake.Piecemove.Y].King = true
	}

	return newBoard
}

func getLoser(board [][]square, startingColor string) string {
	var moves = generateAllValidMoves(board, startingColor)
	if len(moves) > 0 {
		var newMove = randomMove(board, startingColor)
		var newBoard = playmove(newMove, board)
		var newColor = opposingColor(startingColor)
		var newMoves = generateAllValidMoves(newBoard, startingColor)

		if newMove.Piecemove.Flag == "capture" && (len(newMoves) > 0) &&
			newMoves[0].Piecemove.Flag == "capture" {

			newColor = startingColor

		}

		return getLoser(newBoard, newColor)

	}

	return startingColor
}

func backPropagate(node *node, win bool) {

	if win {
		node.wins++
	}

	node.sims++

	if node.parent != nil {

		backPropagate(node.parent, win)
	}

}

func playOut(node *node, rootcolor string) {

	backPropagate(node, rootcolor != getLoser(copyBoard(node.boardState), rootcolor))

}

func expandNode(node *node) {

	var newColor = opposingColor(node.movecolor)

	if node.move.Piecemove.Flag == "capture" {

		newmoves := generateValidMoves(node.boardState, node.move.Piece.Xpos, node.move.Piece.Ypos, node.movecolor)

		if len(newmoves) > 0 && newmoves[0].Flag == "capture" {
			newColor = node.movecolor
		}
	}

	var moves = generateAllValidMoves(node.boardState, newColor)

	for _, move := range moves {

		var newNode = createNewNode(move, node.boardState, newColor, node)
		node.children = append(node.children, newNode)

	}

}

func uctValue(node node) float64 {
	var c = 1.141
	var parentSims int

	if node.parent != nil {
		parentSims = node.parent.sims
	} else {
		parentSims = 1
	}

	if node.sims == 0 {
		return math.MaxInt32
	}

	return ((float64)(node.wins) / (float64)(node.sims)) + (c * math.Pow(math.Log((float64)(parentSims))/(float64)(node.sims), .5))

}

func inRange(x int, y int) bool {
	return (x < 8 && x > -1 && y < 8 && y > -1)
}

func opposingColor(color string) string {
	var str string
	if color == "r" {
		str = "w"
	} else {
		str = "r"
	}
	return str
}

func randomMove(board [][]square, color string) piecemove {
	var moves = generateAllValidMoves(board, color)
	return moves[rand.Intn(len(moves))]
}

func generateAllValidMoves(board [][]square, color string) []piecemove {

	var moves []piecemove
	var captures []piecemove

	for y, row := range board {
		for x := range row {
			if board[x][y].Color == color {
				for _, move := range generateValidMoves(board, x, y, color) {
					if move.Flag == "capture" {
						captures = append(captures, piecemove{piece{x, y}, move})
					} else {
						moves = append(moves, piecemove{piece{x, y}, move})
					}
				}
			}

		}
	}

	if len(captures) > 0 {
		return captures
	}
	return moves

}

func generateValidMoves(board [][]square, x int, y int, color string) []move {

	var moves []move

	var ydirs = []int{-1, 1}

	for _, ydir := range ydirs {
		var xdirs []int
		if board[x][y].King {
			xdirs = []int{1, -1}
		} else if board[x][y].Color == "r" {
			xdirs = append(xdirs, 1)
		} else if board[x][y].Color == "w" {
			xdirs = append(xdirs, -1)
		}

		for _, xdir := range xdirs {

			if inRange(x+xdir, y+ydir) {

				var newMove move
				if board[x+xdir][y+ydir].Color == opposingColor(color) && inRange(x+xdir*2, y+ydir*2) && board[x+xdir*2][y+ydir*2].Color == "null" {
					newMove = move{x + xdir*2, y + ydir*2, "capture", opponentPos{x + xdir, y + ydir}}
					moves = append(moves, newMove)
				} else if board[x+xdir][y+ydir].Color == "null" {
					newMove = move{x + xdir, y + ydir, "move", opponentPos{x + xdir, y + ydir}}
					moves = append(moves, newMove)
				}
			}
		}

	}

	return moves

}

func copyBoard(board [][]square) [][]square {
	newBoard := make([][]square, 8)
	for x := range board {
		for y := range board {
			sq := board[x][y]
			newBoard[x] = append(newBoard[x], sq)
		}
	}

	return newBoard
}

func printBoard(board [][]square) {
	for _, row := range board {
		for _, square := range row {
			fmt.Print("| ")

			if square.King {
				if square.Color == "r" {
					fmt.Print("R")
				} else if square.Color == "w" {
					fmt.Print("W")
				}
			} else {
				fmt.Print(square.Color)
			}

			fmt.Print(" |")
		}
		fmt.Println()
	}
}
