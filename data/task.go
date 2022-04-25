package data

type Task struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
