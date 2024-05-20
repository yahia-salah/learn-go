package repository

import (
	"time"

	"github.com/yahia-salah/learn-go/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestriction) error

	SearchAvailabilityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error)

	SearchAvailabilityForAllRooms(start time.Time, end time.Time) ([]models.Room, error)

	GetRoomByID(id int) (models.Room, error)
}
