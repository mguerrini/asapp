package memory

type memoryHealthRepository struct {
}

func NewMemoryHealthRepositoryRepository () *memoryHealthRepository {
	return &memoryHealthRepository{
	}
}

func (h *memoryHealthRepository) Ping() error {
	return nil
}
