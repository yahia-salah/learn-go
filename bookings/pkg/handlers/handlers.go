package handlers

import (
	"net/http"

	"github.com/yahia-salah/learn-go/pkg/config"
	"github.com/yahia-salah/learn-go/pkg/models"
	"github.com/yahia-salah/learn-go/pkg/render"
)

// The repository used by the handlers
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
}

// Creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// Sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// The Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	stringMap := make(map[string]string)
	stringMap["title"] = "Home"

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap := make(map[string]string)
	stringMap["title"] = "About"
	stringMap["test"] = "Hello, again!"
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Major's Suite page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Major's Suite"

	render.RenderTemplate(w, "majors.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The General's Quarters page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Generals's Quarters"

	render.RenderTemplate(w, "generals.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Contact"

	render.RenderTemplate(w, "contact.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Search for Availability"

	render.RenderTemplate(w, "reservation.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Make Reservation page handler
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Make Reservation"

	render.RenderTemplate(w, "make-reservation.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
