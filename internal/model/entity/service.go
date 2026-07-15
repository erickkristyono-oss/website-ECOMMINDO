package entity

type Service struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ShortDesc   string  `json:"short_desc"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Icon        string  `json:"icon"`
}
