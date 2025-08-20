package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/Vmarcelo49/CardGame/pkg/assets"
	"github.com/Vmarcelo49/CardGame/pkg/card"
	"github.com/Vmarcelo49/CardGame/pkg/logic"
	"github.com/Vmarcelo49/CardGame/pkg/render"
	"github.com/Vmarcelo49/CardGame/pkg/state"
	"github.com/Vmarcelo49/CardGame/pkg/ui"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
)

var (
	BackgroundColor = color.RGBA{R: 31, G: 31, B: 31, A: 255}
	ExitFlag        error
)

type Scene uint8

const (
	MainMenu Scene = iota
	RockPaperScissors
	DuelScene
)

// Game represents the main game state
type Game struct {
	currentScene      Scene
	mainMenuButtons   []*ui.Button
	mouse             *ui.Mouse
	duelRenderer      *render.DuelRenderer
	gamestate         *state.Gamestate
	previousGamestate *state.Gamestate
	exitingDuel       bool
	otherImgs         []*ui.Label
	duelButtons       []*ui.Button
	font              *text.GoTextFaceSource
	inputHandler      *ui.InputHandler
	gameLogic         *logic.GameLogic
}

// NewGame creates a new game instance
func NewGame() (*Game, error) {
	font, err := assets.LoadFont()
	if err != nil {
		return nil, err
	}

	game := &Game{
		mouse:        &ui.Mouse{},
		currentScene: MainMenu,
		font:         font,
		inputHandler: ui.NewInputHandler(),
		gameLogic:    logic.NewGameLogic(),
	}

	return game, nil
}

// Update updates the game state
func (g *Game) Update() error {
	g.mouse.UpdateMouseState()

	switch g.currentScene {
	case DuelScene:
		g.updateGameLogic()
	case MainMenu:
		if g.exitingDuel {
			g.duelRenderer.FreeImages()
			g.duelRenderer = nil
			g.exitingDuel = false
		}
		if g.mainMenuButtons == nil {
			g.mainMenuButtons = g.createMainMenuButtons()
		}
		for _, b := range g.mainMenuButtons {
			if !b.AlreadyClicked {
				ExitFlag = b.CheckClicked(g.mouse)
			}
		}
	case RockPaperScissors:
		fmt.Println("Rock Paper Scissors - Coming soon...")
	}

	return ExitFlag
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.currentScene {
	case DuelScene:
		g.drawDuel(screen)
	case MainMenu:
		render.DrawMainMenu(screen, g.mainMenuButtons)
	case RockPaperScissors:
		fmt.Println("Rock Paper Scissors will be implemented someday...")
	}
}

// Layout returns the screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) updateGameLogic() {
	drawCardFunc := func(target string, amount uint) {
		g.drawCard(target, amount)
	}
	exitDuelFunc := func() {
		g.exitingDuel = true
		g.currentScene = MainMenu
	}

	newGamestate, err := logic.UpdateGameLogic(g.gamestate, g.previousGamestate, g.mouse,
		g.duelButtons, g.inputHandler, drawCardFunc, exitDuelFunc)
	if err == nil {
		g.previousGamestate = newGamestate
	}
}

func (g *Game) drawCard(target string, amount uint) {
	if target == "player" {
		for i := 0; i < int(amount); i++ {
			g.gamestate.P1.DrawCard()
		}
	} else {
		for i := 0; i < int(amount); i++ {
			g.gamestate.P2.DrawCard()
		}
	}
}

func (g *Game) drawDuel(screen *ebiten.Image) {
	if g.duelRenderer != nil && g.gamestate != nil {
		g.duelRenderer.DrawDuel(screen,
			g.gamestate.P1.Hand, g.gamestate.P1.Deck, g.gamestate.Field.P1,
			g.gamestate.P2.Hand, g.gamestate.P2.Deck, g.gamestate.Field.P2,
			g.gamestate.P1.HP, g.gamestate.P2.HP,
			g.duelButtons, g.otherImgs)
	}
}

func (g *Game) createMainMenuButtons() []*ui.Button {
	duelFunc := func() error {
		return g.loadDuelMode()
	}
	deckEditorFunc := func() error {
		fmt.Println("Deck Editor - Coming soon...")
		return nil
	}
	exitFunc := func() error {
		return ebiten.Termination
	}

	buttons := ui.CreateMainMenuButtons(ScreenWidth, ScreenHeight, g.font, duelFunc, deckEditorFunc, exitFunc)
	return buttons
}

func (g *Game) loadDuelMode() error {
	g.currentScene = DuelScene

	renderer, err := render.LoadDuelRenderer()
	if err != nil {
		return fmt.Errorf("failed to load duel renderer: %w", err)
	}
	g.duelRenderer = renderer

	deck := "./deck/testDeck.txt"
	gameState, err := state.NewGameState(deck, deck, assets.NewDeck)
	if err != nil {
		return fmt.Errorf("failed to create game state: %w", err)
	}
	g.gamestate = gameState

	if err := g.setupDuelUI(); err != nil {
		return fmt.Errorf("failed to setup duel UI: %w", err)
	}

	// Reset main menu buttons to avoid being clicked again
	g.mainMenuButtons = nil
	return nil
}

func (g *Game) setupDuelUI() error {
	// Create duel buttons
	g.createDuelButtons()

	// Create labels for HP display
	p1HPLabel := ui.NewTextLabel(fmt.Sprintf("P1 HP: %d", g.gamestate.P1.HP), 50, 300, g.font)
	p2HPLabel := ui.NewTextLabel(fmt.Sprintf("P2 HP: %d", g.gamestate.P2.HP), 50, 150, g.font)

	g.otherImgs = []*ui.Label{p1HPLabel, p2HPLabel}
	return nil
}

func (g *Game) createDuelButtons() {
	// Create end turn button
	endTurnFunc := func() error {
		fmt.Println("End turn clicked")
		return nil
	}

	endTurnButton := ui.NewButton(assets.TurnButtonSize, assets.TurnButtonSize,
		ScreenWidth-assets.TurnButtonSize-20, (ScreenHeight/2)-(assets.TurnButtonSize/2),
		"End Turn", endTurnFunc, g.font)

	// Create card action button (will be positioned dynamically)
	cardActionFunc := func() error {
		fmt.Println("Card action clicked")
		return nil
	}

	cardActionButton := ui.NewButton(80, 30, -5000, -5000, "Use", cardActionFunc, g.font)

	g.duelButtons = []*ui.Button{endTurnButton, cardActionButton}
}
