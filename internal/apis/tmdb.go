package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	apiurl = "https://api.themoviedb.org/3/"
	token  = ""
)

func Init(apitoken string) {
	token = apitoken
}

func request(requrl string) (*http.Response, []byte) {
	req, _ := http.NewRequest("GET", apiurl+requrl, nil)

	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)

	res.Body.Close()
	return res, body
}

func buildQuery(opts map[string]string) string {
	query := url.Values{}
	for k, v := range opts {
		query.Add(k, v)
	}
	return query.Encode()
}

type Response struct {
	Results []Movie `json:"results"`
}

type Movie struct {
	Id            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	PosterPath    string  `json:"poster_path"`
	ReleaseDate   string  `json:"release_date"`
	ImdbID        string  `json:"imdb_id"`
	Lenght        int     `json:"runtime"`
	Genres        []Genre `json:"genres"`
	Overview      string  `json:"overview"`
	Director      string
}

type Genre struct {
	Name string `json:"name"`
}

type Crew struct {
	Person []Person `json:"crew"`
}

type Person struct {
	Name       string `json:"name"`
	Department string `json:"job"`
}

type WatchProviders struct {
	ID      int                `json:"id"`
	Results map[string]Country `json:"results"`
}

type Country struct {
	Link     string     `json:"link"`
	Flatrate []Provider `json:"flatrate,omitempty"`
	Free     []Provider `json:"free,omitempty"`
	Rent     []Provider `json:"rent,omitempty"`
	Buy      []Provider `json:"buy,omitempty"`
}

type Provider struct {
	LogoPath        string `json:"logo_path"`
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
	DisplayPriority int    `json:"display_priority"`
}

func SearchTmdbMovie(search string) []Movie {
	opts := map[string]string{"query": search, "page": "1"}
	url := fmt.Sprintf("search/movie?%s", buildQuery(opts))

	_, body := request(url)

	var responseObjekt Response
	json.Unmarshal(body, &responseObjekt)

	return responseObjekt.Results
}

func getMovieDirector(id string) string {
	url := fmt.Sprintf("movie/%s/credits", id)
	_, body := request(url)

	var crewObjekt Crew
	json.Unmarshal(body, &crewObjekt)

	for _, v := range crewObjekt.Person {
		if v.Department == "Director" {
			return v.Name
		}
	}
	return ""
}

func GetMovieStreamingServices(id int, lang string) Country {
	url := fmt.Sprintf("movie/%d/watch/providers", id)
	_, body := request(url)

	var responseObjekt WatchProviders
	json.Unmarshal(body, &responseObjekt)

	// Returning the available streaming services
	return responseObjekt.Results[lang]
}

func GetTmdbMovie(id string) Movie {
	url := fmt.Sprintf("movie/%s", id)
	_, body := request(url)

	getMovieDirector(id)

	var responseObjekt Movie
	json.Unmarshal(body, &responseObjekt)

	responseObjekt.Director = getMovieDirector(id)

	return responseObjekt
}

type WatchResults struct {
	Results []Provider `json:"results,omitempty"`
}

func GetMovieProviders(region string) []Provider {
	opts := map[string]string{"language": "en_US", "watch_region": region}
	url := fmt.Sprintf("watch/providers/movie?%s", buildQuery(opts))

	_, body := request(url)

	var responseObjekt WatchResults
	json.Unmarshal(body, &responseObjekt)

	return responseObjekt.Results
}

type RegionResults struct {
	Results []Region `json:"results"`
}

type Region struct {
	Iso         string `json:"iso_3166_1"`
	EnglishName string `json:"english_name"`
	NativeName  string `json:"native_name"`
}

func GetAvailableRegions() []Region {
	url := "watch/providers/regions"

	_, body := request(url)

	var responseObjekt RegionResults
	json.Unmarshal(body, &responseObjekt)

	return responseObjekt.Results
}
