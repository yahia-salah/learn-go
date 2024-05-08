package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yahia-salah/learn-go/internal/config"
	"github.com/yahia-salah/learn-go/internal/forms"
	"github.com/yahia-salah/learn-go/internal/models"
	"github.com/yahia-salah/learn-go/internal/render"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
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

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Major's Suite page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Major's Suite"

	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The General's Quarters page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Generals's Quarters"

	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Contact"

	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Search Availability page handler
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Search for Availability"

	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The PostAvailability handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// The Availability JSON handler
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	resp := jsonResponse{OK: true, Message: fmt.Sprintf("There are rooms available from %s to %s", start, end)}

	out, err := json.MarshalIndent(resp, "", "\t")

	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// The Make Reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Make Reservation"

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

// PostReservation handler for posting the reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{FirstName: r.Form.Get("firstName"), LastName: r.Form.Get("lastName"), Email: r.Form.Get("email"), Phone: r.Form.Get("phone")}

	form := forms.New(r.PostForm)

	form.Has("firstName", r)
	form.Has("lastName", r)
	form.Has("email", r)
	form.Has("phone", r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		stringMap := make(map[string]string)
		stringMap["title"] = "Make Reservation"

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			StringMap: stringMap,
			Form:      form,
			Data:      data,
		})

		return
	}
}
