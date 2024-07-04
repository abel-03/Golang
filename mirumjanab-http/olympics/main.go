// +build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

var dataPath string
var athletes []AthleteStruct

const (
	OK = 200
	ERROR = 400
	NotFound = 404
)


func parseJsonData() {
	data, err := ioutil.ReadFile(dataPath)
	if err == nil {
		json.Unmarshal(data, &athletes)	
	}

}

func findName(w *http.ResponseWriter, name string, sport string) AthleteInfo {
	var filteredAthletes []AthleteStruct
	var curr AthleteInfo

	currConutry := ""
	
	i := 0
	for _, val := range athletes {
		if val.Athlete == name {
			if sport != "" && val.Sport != sport {
				continue
			}
			if i != 0 &&  val.Country == currConutry {
				filteredAthletes = append(filteredAthletes, val)
			} else if i == 0 {
				currConutry = val.Country
				filteredAthletes = append(filteredAthletes, val)
			}
			i++
		}
	}
	if len(filteredAthletes) == 0 {
		(*w).WriteHeader(NotFound)
		return AthleteInfo{}
	}
	return add(filteredAthletes, curr, name, currConutry)

}

func add(filteredAthletes []AthleteStruct, curr AthleteInfo, name string, currConutry string) AthleteInfo {
	curr.Athlete = name
	curr.Country = currConutry
	var medalsByYearMap = make(map[int]MedalsAtYear)
	for _, val := range filteredAthletes {
		medalsByYearMap[val.Year] = MedalsAtYear{
			Gold:   medalsByYearMap[val.Year].Gold + val.Gold,
			Silver: medalsByYearMap[val.Year].Silver + val.Silver,
			Bronze: medalsByYearMap[val.Year].Bronze + val.Bronze,
			Total:  medalsByYearMap[val.Year].Total + val.Total,
		}
		curr.Medals.Gold += val.Gold
		curr.Medals.Silver += val.Silver
		curr.Medals.Bronze += val.Bronze
		curr.Medals.Total += val.Total
	}
	curr.MedalsByYears = medalsByYearMap
	return curr
}

func topSport(w *http.ResponseWriter, currSport string, limit int) []AthleteInfo {
	used := make(map[string]bool)
	currAth := make([]AthleteInfo, 0)
	for _, val := range athletes {
		if val.Sport != currSport || used[val.Athlete] {
			continue
		}
		currAth = append(currAth, findName(w, val.Athlete, val.Sport))
		used[val.Athlete] = true
	}
	if len(currAth) == 0 {
		(*w).WriteHeader(NotFound)
		return make([]AthleteInfo, 0)
	}

	return sortAth(currAth, limit)
}

func sortAth(currAth []AthleteInfo, limit int) []AthleteInfo {
	sort.Slice(currAth, func(i, j int) bool {
		if currAth[i].Medals.Gold == currAth[j].Medals.Gold {
			if currAth[i].Medals.Silver == currAth[j].Medals.Silver {
				if currAth[i].Medals.Bronze == currAth[j].Medals.Bronze {
					return currAth[j].Athlete > currAth[i].Athlete
				}
				return currAth[j].Medals.Bronze < currAth[i].Medals.Bronze
			}
			return currAth[j].Medals.Silver < currAth[i].Medals.Silver
		}
		return currAth[j].Medals.Gold < currAth[i].Medals.Gold
	})
	if limit > len(currAth){
		return currAth
	}
	return currAth[:limit]
}

func topCounties(w *http.ResponseWriter, year int, limit int) []TopCountriesInYear {
	countriesMap := make(map[string]TopCountriesInYear)
	for _, val := range athletes {
		if val.Year != year {
			continue
		}
		countriesMap[val.Country] = TopCountriesInYear{
			Country: val.Country,
			Gold:    countriesMap[val.Country].Gold + val.Gold,
			Silver:  countriesMap[val.Country].Silver + val.Silver,
			Bronze:  countriesMap[val.Country].Bronze + val.Bronze,
			Total:   countriesMap[val.Country].Total + val.Total,
		}
	
	}
	if len(countriesMap) == 0 {
		(*w).WriteHeader(NotFound)
		return make([]TopCountriesInYear, 0)
	}
	return sorted(countriesMap, limit)
}

func sorted(countriesMap map[string]TopCountriesInYear, limit int, )  []TopCountriesInYear{
	sortedCountries := make([]TopCountriesInYear, 0)
	countryUsed := make(map[string]bool)
	for _, val := range athletes {
		if countryUsed[val.Country] {
			continue
		}
		sortedCountries = append(sortedCountries, countriesMap[val.Country])
		countryUsed[val.Country] = true
	}
	sort.Slice(sortedCountries, func(i, j int) bool {
		if sortedCountries[i].Gold == sortedCountries[j].Gold {
			if sortedCountries[i].Silver == sortedCountries[j].Silver {
				if sortedCountries[i].Bronze == sortedCountries[j].Bronze {
					return sortedCountries[j].Country > sortedCountries[i].Country
				}
				return sortedCountries[j].Bronze < sortedCountries[i].Bronze
			}
			return sortedCountries[j].Silver < sortedCountries[i].Silver
		}
		return sortedCountries[j].Gold < sortedCountries[i].Gold
	})
	sortedNonNull := make([]TopCountriesInYear, 0)
	for _, val := range sortedCountries {
		if val.Country == "" {
			continue
		}
		sortedNonNull = append(sortedNonNull, val)
	}
	if limit > len(sortedNonNull) {
		return sortedNonNull
	}
	return sortedNonNull[:limit]
}

func handlerName(w http.ResponseWriter, que url.Values) {
	curr := findName(&w, que["name"][0], "")
		if curr.Athlete != "" {
			currentInfoJSON, _ := json.Marshal(curr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(OK)
			fmt.Fprintln(w, string(currentInfoJSON))
		}
}

func handlerSport(w http.ResponseWriter, que url.Values, limit int) {
	var err error
	if len(que["limit"]) >= 1 {
		limit, err = strconv.Atoi(que["limit"][0])
		if err != nil {
			w.WriteHeader(ERROR)
		}
	}
	currentTop := topSport(&w, que["sport"][0], limit)
	if len(currentTop) >= 1 {
		currentTopJSON, _ := json.Marshal(currentTop)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(OK)
		fmt.Fprintln(w, string(currentTopJSON))
	}
}

func handlerYear(w http.ResponseWriter, que url.Values) {
	limit := 3
	var err error
	if len(que["limit"]) >= 1 {
		limit, err = strconv.Atoi(que["limit"][0])
		if err != nil {
			w.WriteHeader(ERROR)
		}
	}
	year, _ := strconv.Atoi(que["year"][0])
	currentTop := topCounties(&w, year, limit)
	if len(currentTop) >= 1 {
		currentTopJSON, _ := json.Marshal(currentTop)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(OK)
		fmt.Fprintln(w, string(currentTopJSON))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := "http://" + r.Host + r.URL.String()
	u, _ := url.Parse(urlPath)
	limit := 3	
	que := u.Query()
	parseJsonData()
	if len(que["name"]) >= 1 {
			handlerName(w, que)
	} else if len(que["sport"]) >= 1 {
		handlerSport(w, que, limit)
	} else if len(que["year"]) >= 1 {
		handlerYear(w, que)
	}
}

func main() {
	portPtr := flag.Int("port", 8000, "port string")
	dataPtr := flag.String("data", "", "data string")
	flag.Parse()
	portNumber := *portPtr
	dataPath = *dataPtr
	http.HandleFunc("/", handler)
	localAddress := "localhost:" + strconv.Itoa(portNumber)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}