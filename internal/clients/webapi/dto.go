package webapi

type carDTO struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	ManufacturerId int      `json:"manufacturerId"`
	CategoryId     int      `json:"categoryId"`
	Year           int      `json:"year"`
	Specs          specsDTO `json:"specifications"`
	Image          string   `json:"image"`
}

type specsDTO struct {
	Engine     string `json:"engine"`
	HP         int    `json:"horsepower"`
	Gearbox    string `json:"transmission"`
	Drivetrain string `json:"drivetrain"`
}

type manufacturerDTO struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type categoryDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
