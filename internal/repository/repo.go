package repository

type Repository struct {
	data map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[string]string),
	}
}
