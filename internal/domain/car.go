package domain

type Car struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ManufacturerId int    `json:"manufacturerId"`
	CategoryId     int    `json:"categoryId"`
	Year           int    `json:"year"`
	Specs          Specs  `json:"specifications"`
	Image          string `json:"image"`
}

type Specs struct {
	Engine     string `json:"engine"`
	HP         int    `json:"horsepower"`
	Gearbox    string `json:"transmission"`
	Drivetrain string `json:"drivetrain"`
}

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CarSummary struct {
	ID      int
	Heading string
	Year    int
	Gearbox string
	Image   string
}

// Example
// "carModels": [
//         {
//         "id": 1,
//         "name": "Toyota Corolla",
//         "manufacturerId": 1,
//         "categoryId": 2,
//         "year": 2023,
//         "specifications": {
//             "engine": "1.8L Inline-4",
//             "horsepower": 139,
//             "transmission": "CVT",
//             "drivetrain": "Front-Wheel Drive"
//         },
//         "image": "toyota_corolla.jpg"
//         }
//     ]
