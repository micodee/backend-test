package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Dice represents a die with a top side value.
type Dice struct {
	topSideVal int
}

// Player represents a player with a cup of dice.
type Player struct {
	diceInCup []Dice
	name      string
	position  int
	point     int
}

// Game represents the dice game with multiple players.
type Game struct {
	players               []*Player
	round                 int
	numberOfPlayer        int
	numberOfDicePerPlayer int
	REMOVED_WHEN_DICE_SIX int
	MOVE_WHEN_DICE_ONE    int
}

// GetTopSideVal returns the top side value of the dice.
func (d *Dice) GetTopSideVal() int {
	return d.topSideVal
}

// Roll rolls the dice and sets the top side value to a random number between 1 and 6.
func (d *Dice) Roll() *Dice {
	d.topSideVal = rand.Intn(6) + 1
	return d
}

// SetTopSideVal sets the top side value of the dice.
func (d *Dice) SetTopSideVal(topSideVal int) *Dice {
	d.topSideVal = topSideVal
	return d
}

// Get Dice Player
func (p *Player) GetDiceInCup() []Dice {
	return p.diceInCup
}

// Get Name Player
func (p *Player) GetName() string {
	return p.name
}

// Get Posisition Player
func (p *Player) GetPosition() int {
	return p.position
}

// Add Point
func (p *Player) AddPoint(point int) {
	p.point += point
}

// Get Point
func (p *Player) GetPoint() int {
	return p.point
}

// Mulai melempar dadu
func (p *Player) Play() {
	for i := range p.diceInCup {
		p.diceInCup[i].Roll()
	}
}

// menampilkan putaran saat ini
func (g *Game) DisplayRound() {
	fmt.Printf("Giliran %d\n", g.round)
}

// Hapus dadu karena angka 6
func (p *Player) RemoveDice(index int) {
	if index >= 0 && index < len(p.diceInCup) {
		p.diceInCup = append(p.diceInCup[:index], p.diceInCup[index+1:]...)
	}
}

// Menambahkan dadu karena angka 1 ke pemain diatasnya
func (p *Player) InsertDice(dice Dice) {
	p.diceInCup = append(p.diceInCup, dice)
}

// Membuat pemain baru dengan jumlah dadu, posisi, dan nama opsional tertentu.
func NewPlayer(numberOfDice int, position int, name string) *Player {
	point := 0
	diceInCup := make([]Dice, numberOfDice)
	return &Player{
		diceInCup: diceInCup,
		name:      name,
		position:  position,
		point:     point,
	}
}

// Membuat permainan baru
func NewGame(numberOfPlayer int, numberOfDicePerPlayer int) *Game {
	round := 0
	players := make([]*Player, numberOfPlayer)

	for i := 0; i < numberOfPlayer; i++ {
		players[i] = NewPlayer(numberOfDicePerPlayer, i, fmt.Sprintf("Pemain #%d", i+1))
	}

	return &Game{
		players:               players,
		round:                 round,
		numberOfPlayer:        numberOfPlayer,
		numberOfDicePerPlayer: numberOfDicePerPlayer,
		REMOVED_WHEN_DICE_SIX: 6,
		MOVE_WHEN_DICE_ONE:    1,
	}
}

// menampilkan nilai dadu & point
func (g *Game) DisplayTopSideDice(title string) {
	fmt.Printf("%s:\n", title)
	for _, player := range g.players {
		if len(player.diceInCup) > 0 {
			fmt.Printf("%s (%d): ", player.name, player.point)
			var diceTopSide string
			for _, dice := range player.diceInCup {
				diceTopSide += fmt.Sprintf("%d, ", dice.topSideVal)
			}
			fmt.Println(diceTopSide[:len(diceTopSide)-2]) // Hapus koma
		} else {
			fmt.Printf("%s (%d): _ (Berhenti bermain karena tidak memiliki dadu)\n", player.name, player.point)
		}
	}
}

// Start the game.
func (g *Game) Start() {
	fmt.Printf("Pemain = %d, Dadu = %d\n\n", g.numberOfPlayer, g.numberOfDicePerPlayer)
	fmt.Println("===============================")

	for {
		g.round++
		diceCarryForward := make(map[int][]Dice)

		for _, player := range g.players {
			player.Play()
		}

		g.DisplayRound()
		g.DisplayTopSideDice("Lempar Dadu")

		for index, player := range g.players {
			tempDiceArray := make([]Dice, 0)

			for diceIndex, dice := range player.diceInCup {
				if dice.topSideVal == g.REMOVED_WHEN_DICE_SIX {
					player.AddPoint(1)
					player.RemoveDice(diceIndex)
				}

				if dice.topSideVal == g.MOVE_WHEN_DICE_ONE {
					nextPlayerIndex := (index + 1) % g.numberOfPlayer
					if len(g.players[nextPlayerIndex].diceInCup) > 0 {
						g.players[nextPlayerIndex].InsertDice(dice)
						player.RemoveDice(diceIndex)
						break // Keluar dari loop setelah memindahkan dadu
					} else {
						tempDiceArray = append(tempDiceArray, dice)
						player.RemoveDice(diceIndex)
					}
				}
			}

			diceCarryForward[(index+1)%g.numberOfPlayer] = append(diceCarryForward[(index+1)%g.numberOfPlayer], tempDiceArray...)
		}

		g.DisplayTopSideDice("Setelah Evaluasi")

		playerHasDice := g.numberOfPlayer

		for _, player := range g.players {
			if len(player.diceInCup) <= 0 {
				playerHasDice--
			}
		}

		if playerHasDice <= 1 {
			winner := g.GetWinner()
			g.DisplayWinner(winner)
			break
		}
		fmt.Println("===============================")
	}
}

// Get Winner dengan pemain jumlah point tertinggi
func (g *Game) GetWinner() *Player {
	var winner *Player
	highscore := 0
	for _, player := range g.players {
		if player.point > highscore {
			highscore = player.point
			winner = player
		}
	}
	return winner
}

// Display the Winner
func (g *Game) DisplayWinner(player *Player) {
	fmt.Println()
	fmt.Printf("Game dimenangkan oleh %s karena memiliki poin lebih banyak dari pemain lainnya", player.name)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var pemain, dadu int
	fmt.Print("Masukkan jumlah pemain: ")
	fmt.Scan(&pemain)
	fmt.Print("Masukkan jumlah dadu: ")
	fmt.Scan(&dadu)

	game := NewGame(pemain, dadu)
	game.Start()
}
