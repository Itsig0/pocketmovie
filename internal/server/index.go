package server

import (
	"slices"

	"git.itsigo.dev/istigo/pocketmovie/cmd/web"
	"git.itsigo.dev/istigo/pocketmovie/internal/apis"
	"github.com/gofiber/fiber/v3"
)

func (s *FiberServer) Index(c fiber.Ctx) error {
	movies, _ := s.db.ListMovies(c)
	settings, _ := s.db.ListSettings(c)
	return render(c, web.Show(movies, settings))
}

func (s *FiberServer) Watchlist(c fiber.Ctx) error {
	movies, _ := s.db.ListWatchlist(c)
	return render(c, web.WatchList(movies))
}

func (s *FiberServer) Settings(c fiber.Ctx) error {
	ls, _ := s.db.ListSreamingServices(c)
	_, ts := s.getAllStreamingServices(c)
	slices.Sort(ts)

	reg := apis.GetAvailableRegions()
	selectedreg, _ := s.db.ListSetting(c, 3)
	key, _ := s.db.ListSetting(c, 2)

	config := web.SettingsConfig{
		Providers:          ls,
		AvailableProviders: ts,
		Regions:            reg,
		SelectedRegion:     selectedreg.Value,
		APIKey:             key.Value,
	}

	return render(c, web.Settings(config))
}
