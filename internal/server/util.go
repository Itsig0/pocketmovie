package server

import (
	"slices"
	"strings"

	"git.itsigo.dev/istigo/pocketmovie/internal/apis"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var store = session.New()

func render(c fiber.Ctx, component templ.Component) error {
	// or templ wil bork...
	c.Set("Content-type", "text/html")
	return component.Render(c, c.Response().BodyWriter())
}

// Gets the Streaming services for a movie as a string (service 1, service 2, ...)
func (s *FiberServer) getStreamingServicesForMovie(c fiber.Ctx, id int64) string {
	// Get the streaming providers
	dbServicesSlice := s.getLocalStreamingServices(c)

	region, _ := s.db.ListSetting(c, 3)

	flatrate := apis.GetMovieStreamingServices(int(id), region.Value).Flatrate
	free := apis.GetMovieStreamingServices(int(id), region.Value).Free
	flatrate = append(flatrate, free...)
	streaningservices := []string{}
	for _, v := range flatrate {
		if slices.Contains(dbServicesSlice, v.ProviderName) {
			streaningservices = append(streaningservices, v.ProviderName)
		}
	}
	return strings.Join(streaningservices, ", ")
}

func (s *FiberServer) getAllStreamingServices(c fiber.Ctx) ([]string, []string) {
	ls := s.getLocalStreamingServices(c)
	ts := s.getTmdbStreamingServices(c, ls)
	return ls, ts
}

func (s *FiberServer) getLocalStreamingServices(c fiber.Ctx) []string {
	services, _ := s.db.ListSreamingServices(c)

	sl := []string{}
	for _, p := range services {
		sl = append(sl, p.Title)
	}

	return sl
}

func (s *FiberServer) getTmdbStreamingServices(c fiber.Ctx, localservices []string) []string {
	region, _ := s.db.ListSetting(c, 3)
	services := apis.GetMovieProviders(region.Value)

	sl := []string{}

	for _, p := range services {
		if !slices.Contains(localservices, p.ProviderName) {
			sl = append(sl, p.ProviderName)
		}
	}

	return sl
}
