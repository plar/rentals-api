package domain

type Rental struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Type            string   `json:"type"`
	Make            string   `json:"make"`
	Model           string   `json:"model"`
	Year            int      `json:"year"`
	Length          float64  `json:"length"`
	Sleeps          int      `json:"sleeps"`
	PrimaryImageURL string   `json:"primary_image_url"`
	Price           Price    `json:"price"`
	Location        Location `json:"location"`
	User            User     `json:"user"`
}

type Price struct {
	Day int `json:"day"`
}

type Location struct {
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zip     string  `json:"zip"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
