package dbrepo

import (
	"context"
	"time"

	"github.com/yahia-salah/learn-go/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation in the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	statement := `insert into reservations (first_name, last_name, email, phone, start_date, end_date,
		room_id, created_at, updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, statement,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StarDate,
		res.EndDate,
		res.RoomID, time.Now(), time.Now()).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// Inserts a room restriction in the database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `insert into room_restrictions (start_date, end_date,
		room_id, reservation_id, restriction_id, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, statement,
		res.StarDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionID,
		time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

// Checks if a room is available for reservation
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `select count(id) 
	from room_restrictions
	where
	$1 <= end_date and $2 >= start_date and room_id=$3;
	`

	row := m.DB.QueryRowContext(ctx, query,
		start, end, roomId)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}

// Checks all rooms' available for reservation
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start time.Time, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select r.id,r.room_name 
	from rooms r
	where r.id not in (
		select room_id from room_restrictions
		where
		$1 <= end_date and $2 >= start_date
	);
	`

	rows, err := m.DB.QueryContext(ctx, query,
		start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	err = rows.Err()
	if err != nil {
		return rooms, err
	}

	return rooms, nil
}

// Gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select id,room_name,created_at,updated_at from rooms where id=$1;`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return room, err
	}
	return room, nil
}
