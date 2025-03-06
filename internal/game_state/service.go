package gameState

// Service handles all game-related operations
type Service struct {
	gameStateRepo GameStateRepository
}

// NewService creates a new game service
func NewService(
	gameStateRepo GameStateRepository,
) *Service {
	return &Service{
		gameStateRepo: gameStateRepo,
	}
}
