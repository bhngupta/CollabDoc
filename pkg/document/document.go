package document

type Document struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Version   int    `json:"version"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
