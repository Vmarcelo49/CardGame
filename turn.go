package main

type Phase uint8

func (p Phase) String() string {
	switch p {
	case DrawPhase:
		return "DrawPhase"
	case MainPhase:
		return "MainPhase"
	case BattlePhase:
		return "BattlePhase"
	case EndPhase:
		return "EndPhase"
	default:
		return "Invalid Phase"
	}
}

const (
	WhoGoesFirst Phase = iota
	DrawPhase
	MainPhase
	BattlePhase
	EndPhase
)

func runPhase(g *Game) {
	switch g.phase {
	case DrawPhase:
		drawPhase(g)
	case MainPhase:
		mainPhase(g)
	case BattlePhase:
		battlePhase(g)
	case EndPhase:
		endPhase(g)
	}
}

func changePhase(g *Game, p Phase) {
	switch p {
	case DrawPhase:
		g.phase = MainPhase
	case MainPhase:
		g.phase = BattlePhase
	case BattlePhase:
		g.phase = EndPhase
	case EndPhase:
		g.phase = DrawPhase
	}
}
