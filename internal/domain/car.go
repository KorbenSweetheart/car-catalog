package domain

const (
	TransmissionManual    = "Manual"
	TransmissionAutomatic = "Automatic"
)

type Car struct {
	ID           int
	Name         string
	Year         int
	Image        string
	Specs        Specs
	Manufacturer Manufacturer
	Category     Category
}

type Specs struct {
	Engine       string
	HP           int
	Gearbox      string
	Transmission string
	Drivetrain   string
}

type Manufacturer struct {
	ID           int
	Name         string
	Country      string
	FoundingYear int
}

type Category struct {
	ID   int
	Name string
}

// used for user input in catalog to filter cars
type FilterOptions struct {
	ManufacturerID int
	CategoryID     int
	MinYear        int
	MinHP          int
	Transmission   string
	Drivetrain     string
}

type Metadata struct {
	Manufacturers []Manufacturer
	Categories    []Category
	Drivetrains   []string // e.g. "All-Wheel Drive", "Rear-Wheel Drive", "Front-Wheel Drive"
	Transmissions []string // e.g. "Automatic", "Manual"
}
