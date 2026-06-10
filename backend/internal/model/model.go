package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. User: โครงสร้างข้อมูลผู้ใช้งาน (Collection: users)
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GoogleID  string             `bson:"google_id" json:"google_id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// 2. Seat: โครงสร้างข้อมูลที่นั่งย่อย (Sub-document อยู่ภายในคอลเลกชัน showtimes)
type Seat struct {
	SeatNo      string     `bson:"seat_no" json:"seat_no"`
	Status      string     `bson:"status" json:"status"` // AVAILABLE, LOCKED, BOOKED
	LockedBy    *string    `bson:"locked_by,omitempty" json:"locked_by,omitempty"`
	LockedUntil *time.Time `bson:"locked_until,omitempty" json:"locked_until,omitempty"`
}

// 3. Movie: โครงสร้างข้อมูลภาพยนตร์ (Collection: movies)
type Movie struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title        string             `bson:"title" json:"title"`
	DurationMins int                `bson:"duration_mins" json:"duration_mins"`
	PosterURL    string             `bson:"poster_url" json:"poster_url"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// 4. Showtime: โครงสร้างข้อมูลรอบฉายและผังที่นั่ง (Collection: showtimes)
type Showtime struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID      primitive.ObjectID `bson:"movie_id" json:"movie_id"`
	StartTime    time.Time          `bson:"start_time" json:"start_time"`
	PricePerSeat float64            `bson:"price_per_seat" json:"price_per_seat"`
	Seats        []Seat             `bson:"seats" json:"seats"`
}

// 5. Booking: โครงสร้างข้อมูลประวัติการจอง (Collection: bookings)
type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShowtimeID primitive.ObjectID `bson:"showtime_id" json:"showtime_id"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Seats      []string           `bson:"seats" json:"seats"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	Status     string             `bson:"status" json:"status"` // PENDING, SUCCESS, TIMEOUT
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

// 6. AuditLog: โครงสร้างข้อมูลบันทึกประวัติเหตุการณ์สำคัญ (Collection: audit_logs)
type AuditLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventType string             `bson:"event_type" json:"event_type"` // BOOKING_SUCCESS, BOOKING_TIMEOUT, SEAT_RELEASED, SYSTEM_ERROR
	Details   string             `bson:"details" json:"details"`
	UserID    string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}