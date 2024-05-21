package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yahia-salah/learn-go/internal/config"
	"github.com/yahia-salah/learn-go/internal/driver"
	"github.com/yahia-salah/learn-go/internal/forms"
	"github.com/yahia-salah/learn-go/internal/helpers"
	"github.com/yahia-salah/learn-go/internal/models"
	"github.com/yahia-salah/learn-go/internal/render"
	"github.com/yahia-salah/learn-go/internal/repository"
	"github.com/yahia-salah/learn-go/internal/repository/dbrepo"
)

// The repository used by the handlers
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// Sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// The Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Home"

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "About"
	stringMap["test"] = "Hello, again!"

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Major's Suite page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Major's Suite"

	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The General's Quarters page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Generals's Quarters"

	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Contact"

	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The Search Availability page handler
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Search for Availability"

	if m.App.Session.Exists(r.Context(), "startDate") && m.App.Session.Exists(r.Context(), "endDate") {
		stringMap["startDate"] = m.App.Session.GetString(r.Context(), "startDate")
		stringMap["endDate"] = m.App.Session.GetString(r.Context(), "endDate")
		m.App.Session.Remove(r.Context(), "startDate")
		m.App.Session.Remove(r.Context(), "endDate")
	}

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// The PostAvailability handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	// 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate,
		err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate,
		err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No rooms available")
		m.App.Session.Put(r.Context(), "startDate", sd)
		m.App.Session.Put(r.Context(), "endDate", ed)
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	res := models.Reservation{
		StarDate: startDate,
		EndDate:  endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	data := make(map[string]interface{})
	data["rooms"] = rooms

	stringMap := make(map[string]string)
	stringMap["title"] = "Choose A Room"

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// The Availability JSON handler
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomID)
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}

	out, err := json.MarshalIndent(resp, "", "\t")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Takes URL parameters, builds session variable, redirects to make reservation page
func (m Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("start_date")
	ed := r.URL.Query().Get("end_date")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	reservation := models.Reservation{
		StarDate: startDate,
		EndDate:  endDate,
		RoomID:   roomId,
	}

	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.Room = room

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Something went wrong"))
		return
	}

	reservation.RoomID = roomId
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// The Make Reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room = room

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StarDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	data := make(map[string]interface{})
	data["reservation"] = res

	stringMap := make(map[string]string)
	stringMap["title"] = "Make Reservation"
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Form:      forms.New(nil),
		Data:      data,
	})
}

// PostReservation handler for posting the reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Something went wrong"))
		return
	}

	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("firstName")
	reservation.LastName = r.Form.Get("lastName")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("firstName", "lastName", "email", "phone")
	form.MinLength("firstName", 5, r)
	form.MinLength("lastName", 5, r)
	form.IsPhone("phone")
	form.IsEmail("email")

	if !form.Valid() {
		sd := reservation.StarDate.Format("2006-01-02")
		ed := reservation.EndDate.Format("2006-01-02")

		data := make(map[string]interface{})
		data["reservation"] = reservation

		stringMap := make(map[string]string)
		stringMap["title"] = "Make Reservation"
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			StringMap: stringMap,
			Form:      form,
			Data:      data,
		})

		return
	}

	reservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StarDate:      reservation.StarDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: reservationId,
		Reservation:   reservation,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// show reservation summary
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Can't cast reservation object from session")
		m.App.Session.Put(r.Context(), "error", "Can't cast reservation object from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	sd := reservation.StarDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	defer m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	stringMap := make(map[string]string)
	stringMap["title"] = "Reservation Summary"
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}
