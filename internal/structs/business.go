package structs

type BusinessConfig struct {
	BusinessList []Business `json:"business"`
}

type Business struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ctime       string `json:"ctime"`
	Status      uint8  `json:"status"`
}
