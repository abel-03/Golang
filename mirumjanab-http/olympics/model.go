package main


type AthleteStruct struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type AthleteInfo struct {
	Athlete       string               `json:"athlete"`
	Country       string               `json:"country"`
	Medals        MedalsAtYear         `json:"medals"`
	MedalsByYears map[int]MedalsAtYear `json:"medals_by_year"`
}

type MedalsAtYear struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type TopCountriesInYear struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}