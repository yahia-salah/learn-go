package models

import "time"

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StarDate  time.Time
	EndDate   time.Time
	RoomID    int
	Room      Room
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomRestriction is the room restriction model
type RoomRestriction struct {
	ID            int
	StarDate      time.Time
	EndDate       time.Time
	RoomID        int
	Room          Room
	ReservationID int
	Reservation   Reservation
	RestrictionID int
	Restriction   Restriction
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
