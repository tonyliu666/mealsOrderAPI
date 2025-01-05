package database

type DBManager interface {
	Save() error
	Read() error
}

func NewDBManager(kind string) DBManager {
	switch kind {
	case "user":
		return &Client{}
	case "diets":
		return &Diets{}
	case "ingredients":
		return &Ingredients{}
	default:
		return nil
	}
}
