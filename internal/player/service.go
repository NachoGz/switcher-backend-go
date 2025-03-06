package player

// Service handles all game-related operations
type Service struct {
	playerRepo PlayerRepository
}

// NewService creates a new game service
func NewService(
	playerRepo PlayerRepository,
) *Service {
	return &Service{
		playerRepo: playerRepo,
	}
}
