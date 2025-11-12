package server

import (
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"git.itsigo.dev/istigo/pocketmovie/cmd/web"
	"git.itsigo.dev/istigo/pocketmovie/internal/apis"
	"git.itsigo.dev/istigo/pocketmovie/internal/database"
	"git.itsigo.dev/istigo/pocketmovie/internal/middleware/apikeychecker"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func (s *FiberServer) RegisterFiberRoutes() {

	s.App.Use(
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed, // 1
		}),
	)

	s.App.Use("/assets", static.New("./assets", static.Config{
		FS:     web.Files,
		Browse: false,
	}))

	s.App.Use("/movie/posters", static.New("./data/img"))

	s.App.Use("/apikey", s.apikeyinput)

	s.App.Use(
		//basicauth.New(basicauth.Config{
		//	Users: map[string]string{
		//		// "doe" hashed using SHA-256
		//		"john": "{SHA256}eZ75KhGvkY4/t0HfQpNPO1aO0tk6wd908bjUGieTKm8=",
		//	},
		//}),
		apikeychecker.New(apikeychecker.Config{DB: *s.db}),
	)

	s.App.Get("/", s.Index)
	s.App.Get("/watchlist", s.Watchlist)
	s.App.Get("/settings", s.Settings)
	s.App.Get("/movie/:id", s.movieDetails)

	s.App.Get("/components/mediaTypeSelect/:id", s.mediaTypeSelect)

	s.App.Post("/apis/tmdb/searchMovie", s.searchMovie)
	s.App.Post("/db/addMovie/:id", s.addMovieToDb)
	s.App.Post("/db/changeMovieStatus/:id.:status", s.changeMovieStatus)
	s.App.Post("/db/changeMovieRating/:id.:rating", s.changeMovieRating)
	s.App.Post("/db/changeMovieOwned/:id", s.changeMovieOwned)
	s.App.Post("/db/changeMovieRipped/:id.:ripped", s.changeMovieRipped)
	s.App.Post("/db/updateStreamingServices", s.updateMovieStreamingServices)
	s.App.Post("/db/updateTableView", s.updateTableView)
	s.App.Post("/db/updateRegion", s.updateRegion)
	s.App.Post("/db/updateApiKey", s.updateApiKey)
	s.App.Delete("/db/deleteMovie/:id", s.deleteMovie)

	s.App.Post("/db/addStreamingService", s.addStreamingService)
	s.App.Delete("/db/deleteStreamingService/:id", s.deleteStreamingService)
}

func (s *FiberServer) apikeyinput(c fiber.Ctx) error {
	key := c.FormValue("apikey")
	apierror := ""
	if key != "" {
		url := "https://api.themoviedb.org/3/authentication"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+key)
		res, _ := http.DefaultClient.Do(req)

		if res.StatusCode == 200 {
			s.db.ChangeSettingValue(c, database.ChangeSettingValueParams{
				ID:    2,
				Value: key,
			})
			c.Redirect().Status(fiber.StatusMovedPermanently).To("/")
		}
		apierror = "Could not certify API Read Access Token"
	}
	return render(c, web.Apikeyinput(apierror))
}

func (s *FiberServer) updateApiKey(c fiber.Ctx) error {
	key := c.FormValue("apikey")
	s.db.ChangeSettingValue(c, database.ChangeSettingValueParams{
		ID:    2,
		Value: key,
	})

	config := web.SettingsConfig{
		APIKey: key,
	}
	return render(c, web.ApiKey(config))
}

func (s *FiberServer) movieDetails(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	movie, _ := s.db.ListMovie(c, id)
	return render(c, web.MovieDetails(movie))
}

func (s *FiberServer) mediaTypeSelect(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	return render(c, web.MovieDetailsOwnedSelect(int8(id)))
}

func (s *FiberServer) searchMovie(c fiber.Ctx) error {
	res := apis.SearchTmdbMovie(c.FormValue("search"))
	return render(c, web.SearchMovie(res))
}

func (s *FiberServer) changeMovieStatus(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	status, _ := strconv.ParseInt(c.Params("status"), 10, 64)

	s.db.ChangeMovieStatus(c, database.ChangeMovieStatusParams{
		ID:     id,
		Status: status,
	})

	return render(c, web.MovieDetailsWatched(id, status))
}

func (s *FiberServer) changeMovieRating(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	rating, _ := strconv.ParseInt(c.Params("rating"), 10, 64)

	s.db.ChangeMovieRating(c, database.ChangeMovieRatingParams{
		ID:     id,
		Rating: rating,
	})

	return render(c, web.MovieDetailsRating(id, rating))
}

func (s *FiberServer) changeMovieOwned(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	option := c.FormValue("option")
	var owned int64 = 0
	if option != "" {
		owned = 1
	}

	s.db.ChangeMovieOwned(c, database.ChangeMovieOwnedParams{
		ID:        id,
		Owned:     owned,
		OwnedType: option,
	})

	return render(c, web.MovieDetailsOwned(id, option))
}

func (s *FiberServer) changeMovieRipped(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	ripped, _ := strconv.ParseInt(c.Params("ripped"), 10, 64)

	s.db.ChangeMovieRipped(c, database.ChangeMovieRippedParams{
		ID:     id,
		Ripped: ripped,
	})

	return render(c, web.MovieDetailsRipped(id, ripped))
}

func (s *FiberServer) updateMovieStreamingServices(c fiber.Ctx) error {
	movies, _ := s.db.ListMovies(c)
	for _, v := range movies {
		s.db.ChangeMovieStreamingServices(c, database.ChangeMovieStreamingServicesParams{
			ID:                v.ID,
			StreamingServices: s.getStreamingServicesForMovie(c, v.Tmdbid),
		})
	}
	movies, _ = s.db.ListWatchlist(c)
	return render(c, web.List(movies))
}

func (s *FiberServer) updateTableView(c fiber.Ctx) error {
	setting, _ := s.db.ListSetting(c, 1)
	v := "false"
	if setting.Value == "false" {
		v = "true"
	}
	s.db.ChangeSettingValue(c, database.ChangeSettingValueParams{
		Value: v,
		ID:    1,
	})
	movies, _ := s.db.ListMovies(c)
	if v == "true" {
		return render(c, web.MovieTiles(movies))
	}
	return render(c, web.MovieList(movies))
}

func (s *FiberServer) updateRegion(c fiber.Ctx) error {
	region := c.FormValue("region")
	s.db.ChangeSettingValue(c, database.ChangeSettingValueParams{
		Value: region,
		ID:    3,
	})

	regions := apis.GetAvailableRegions()

	config := web.SettingsConfig{
		Regions:        regions,
		SelectedRegion: region,
	}
	return render(c, web.Region(config))
}

func (s *FiberServer) addStreamingService(c fiber.Ctx) error {
	title := c.FormValue("service")

	s.db.AddSreamingService(c, title)

	ls, _ := s.db.ListSreamingServices(c)
	_, ts := s.getAllStreamingServices(c)
	slices.Sort(ts)

	config := web.SettingsConfig{
		Providers:          ls,
		AvailableProviders: ts,
	}

	return render(c, web.ProviderTable(config))
}

func (s *FiberServer) deleteStreamingService(c fiber.Ctx) error {
	id := fiber.Params[int64](c, "id")

	s.db.DeleteSreamingService(c, id)

	ls, _ := s.db.ListSreamingServices(c)
	_, ts := s.getAllStreamingServices(c)
	slices.Sort(ts)

	config := web.SettingsConfig{
		Providers:          ls,
		AvailableProviders: ts,
	}

	return render(c, web.ProviderTable(config))
}

func (s *FiberServer) deleteMovie(c fiber.Ctx) error {
	id := fiber.Params[int64](c, "id")

	s.db.DeleteMovie(c, id)

	movies, _ := s.db.ListMovies(c)
	return render(c, web.MovieList(movies))
}

func (s *FiberServer) addMovieToDb(c fiber.Ctx) error {

	res := apis.GetTmdbMovie(c.Params("id"))
	referer := string(c.Request().Header.Referer())

	// Combine Genres
	genres := []string{}
	for _, v := range res.Genres {
		genres = append(genres, v.Name)
	}

	// Get the streaming providers
	streamingservices := s.getStreamingServicesForMovie(c, int64(res.Id))

	rating := 0
	if c.FormValue("rating") != "" {
		rating, _ = strconv.Atoi(c.FormValue("rating"))
	}
	watched := 0
	if c.FormValue("watched") == "on" {
		watched = 1
	}
	watchcount := 0
	if c.FormValue("watchcount") != "" {
		watchcount, _ = strconv.Atoi(c.FormValue("watchcount"))
	}
	owned := 0
	if c.FormValue("owned") == "on" {
		owned = 1
	}
	ownedType := ""
	if c.FormValue("version") != "" {
		ownedType = c.FormValue("version")
	}
	ripped := 0
	if c.FormValue("ripped") == "on" {
		ripped = 1
	}

	movie, err := s.db.CreateMovie(c, database.CreateMovieParams{
		Title:             res.Title,
		OriginalTitle:     res.OriginalTitle,
		Imdbid:            res.ImdbID,
		Tmdbid:            int64(res.Id),
		Length:            int64(res.Lenght),
		Genre:             strings.Join(genres, ", "),
		StreamingServices: streamingservices,
		Director:          res.Director,
		Year:              res.ReleaseDate,
		Watchcount:        int64(watchcount),
		Rating:            int64(rating),
		Status:            int64(watched),
		Owned:             int64(owned),
		OwnedType:         ownedType,
		Ripped:            int64(ripped),
		Review:            "",
		Overview:          res.Overview,
	})
	if err != nil {
		log.Error(err)
	}

	// Get the movie poster
	resp, err := http.Get("https://image.tmdb.org/t/p/original" + res.PosterPath)
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()

	tmp, err := jpeg.Decode(resp.Body)
	if err == nil {
		path := fmt.Sprintf("./data/img/%d.jpg", movie.ID)
		log.Info(path)
		outputFile, err := os.Create(path)
		if err != nil {
			log.Error(err)
		}

		jpeg.Encode(outputFile, tmp, nil)

		outputFile.Close()
	}

	if strings.Contains(referer, "watchlist") {
		movies, _ := s.db.ListWatchlist(c)
		return render(c, web.List(movies))
	}

	gridSetting, _ := s.db.ListSetting(c, 1)
	movies, _ := s.db.ListMovies(c)
	componet := web.MovieList(movies)

	if gridSetting.Value == "true" {
		componet = web.MovieTiles(movies)
	}

	return render(c, componet)
}
