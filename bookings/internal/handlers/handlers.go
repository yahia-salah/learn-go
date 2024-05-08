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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	stringMap := make(map[string]string)
	stringMap["title"] = "Make Reservation"

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Form:      forms.New(nil),
		Data:      data,
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

	form.Required("firstName", "lastName", "email", "phone")
	form.MinLength("firstName", 5, r)
	form.MinLength("lastName", 5, r)
	form.IsPhone("phone")
	form.IsEmail("email")

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

	// show reservation summary
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		log.Println("Can't cast reservation object from session")
		m.App.Session.Put(r.Context(), "error", "Can't cast reservation object from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	defer m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	stringMap := make(map[string]string)
	stringMap["title"] = "Reservation Summary"

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}
