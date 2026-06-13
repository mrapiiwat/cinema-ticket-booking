package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

const lockTTL = 5 * time.Minute

type bookingService struct {
	bookingRepo  ports.BookingRepositoryPort
	showtimeRepo ports.ShowtimeRepositoryPort
	lockRepo     ports.LockRepositoryPort
	mqRepo       ports.MQRepositoryPort
	wsHub        *WSHub
}

func NewBookingService(
	bookingRepo ports.BookingRepositoryPort,
	showtimeRepo ports.ShowtimeRepositoryPort,
	lockRepo ports.LockRepositoryPort,
	mqRepo ports.MQRepositoryPort,
	wsHub *WSHub,
) ports.BookingServicePort {
	return &bookingService{
		bookingRepo:  bookingRepo,
		showtimeRepo: showtimeRepo,
		lockRepo:     lockRepo,
		mqRepo:       mqRepo,
		wsHub:        wsHub,
	}
}

func normalizeSeats(seats []string) []string {
	unique := map[string]bool{}
	for _, seat := range seats {
		value := strings.ToUpper(strings.TrimSpace(seat))
		if value != "" {
			unique[value] = true
		}
	}

	result := make([]string, 0, len(unique))
	for seat := range unique {
		result = append(result, seat)
	}
	sort.Strings(result)
	return result
}

func lockKey(showtimeID string, seatNo string) string {
	return fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeID, seatNo)
}

func randomToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *bookingService) LockSeats(ctx context.Context, showtimeID string, requestedSeats []string, userID string) (*model.Booking, error) {
	seats := normalizeSeats(requestedSeats)
	if len(seats) == 0 {
		return nil, errors.New("at least one seat is required")
	}

	if err := s.showtimeRepo.ReleaseExpiredLocks(ctx, showtimeID, time.Now()); err != nil {
		return nil, err
	}

	showtime, err := s.showtimeRepo.FindByID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}
	if showtime == nil {
		return nil, errors.New("showtime not found")
	}
	if err := validateRequestedSeats(showtime, seats); err != nil {
		return nil, err
	}

	lockToken, err := randomToken()
	if err != nil {
		return nil, err
	}

	lockExpiresAt := time.Now().Add(lockTTL)
	acquired := make([]string, 0, len(seats))
	dbLocked := make([]string, 0, len(seats))

	for _, seatNo := range seats {
		key := lockKey(showtimeID, seatNo)
		ok, err := s.lockRepo.AcquireLock(ctx, key, lockToken, lockTTL)
		if err != nil {
			s.rollbackLocks(ctx, showtimeID, userID, lockToken, acquired, dbLocked)
			s.publishAuditLog(ctx, "SYSTEM_ERROR", fmt.Sprintf("Redis lock failed for seat %s: %v", seatNo, err), userID)
			return nil, err
		}
		if !ok {
			s.rollbackLocks(ctx, showtimeID, userID, lockToken, acquired, dbLocked)
			s.publishAuditLog(ctx, "SYSTEM_ERROR", fmt.Sprintf("Seat %s is already locked", seatNo), userID)
			return nil, errors.New("one or more seats are already locked")
		}
		acquired = append(acquired, seatNo)

		if err := s.showtimeRepo.LockSeatIfAvailable(ctx, showtimeID, seatNo, userID, lockExpiresAt); err != nil {
			s.rollbackLocks(ctx, showtimeID, userID, lockToken, acquired, dbLocked)
			s.publishAuditLog(ctx, "SYSTEM_ERROR", fmt.Sprintf("Mongo seat lock failed for seat %s: %v", seatNo, err), userID)
			return nil, errors.New("one or more seats are not available")
		}
		dbLocked = append(dbLocked, seatNo)
	}

	booking := &model.Booking{
		ShowtimeID:    showtime.ID,
		UserID:        userID,
		Seats:         seats,
		TotalPrice:    showtime.PricePerSeat * float64(len(seats)),
		Status:        model.BookingStatusPending,
		LockToken:     lockToken,
		LockExpiresAt: lockExpiresAt,
	}
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		s.rollbackLocks(ctx, showtimeID, userID, lockToken, acquired, dbLocked)
		s.publishAuditLog(ctx, "SYSTEM_ERROR", fmt.Sprintf("Booking create failed: %v", err), userID)
		return nil, err
	}

	s.publishAuditLog(ctx, "SEAT_LOCKED", fmt.Sprintf("Booking %s locked seats %s", booking.ID.Hex(), strings.Join(seats, ",")), userID)
	s.broadcastSeats(showtimeID, seats, model.SeatStatusLocked, userID, booking.ID.Hex(), lockExpiresAt)
	s.scheduleBookingExpiration(booking.ID.Hex(), time.Until(lockExpiresAt))

	return booking, nil
}

func (s *bookingService) ConfirmBooking(ctx context.Context, bookingID, userID string) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil || booking == nil {
		return nil, errors.New("booking not found")
	}
	if booking.UserID != userID {
		return nil, errors.New("booking belongs to another user")
	}
	if booking.Status != model.BookingStatusPending {
		return nil, errors.New("booking is already processed")
	}
	if time.Now().After(booking.LockExpiresAt) {
		_ = s.expireBooking(ctx, booking)
		return nil, errors.New("booking lock expired")
	}

	showtimeID := booking.ShowtimeID.Hex()
	for _, seatNo := range booking.Seats {
		value, err := s.lockRepo.GetLockValue(ctx, lockKey(showtimeID, seatNo))
		if err != nil {
			return nil, err
		}
		if value != booking.LockToken {
			_ = s.expireBooking(ctx, booking)
			return nil, errors.New("booking lock is no longer active")
		}
	}

	for _, seatNo := range booking.Seats {
		if err := s.showtimeRepo.UpdateSeatStatus(ctx, showtimeID, seatNo, model.SeatStatusBooked, &booking.UserID, nil); err != nil {
			s.publishAuditLog(ctx, "SYSTEM_ERROR", fmt.Sprintf("Booking %s failed to mark seat %s as booked: %v", bookingID, seatNo, err), booking.UserID)
			return nil, err
		}
	}

	if err := s.bookingRepo.UpdateStatus(ctx, bookingID, model.BookingStatusSuccess); err != nil {
		return nil, err
	}
	for _, seatNo := range booking.Seats {
		_ = s.lockRepo.ReleaseLock(ctx, lockKey(showtimeID, seatNo), booking.LockToken)
	}

	booking.Status = model.BookingStatusSuccess
	now := time.Now()
	booking.ConfirmedAt = &now
	s.publishAuditLog(ctx, "BOOKING_SUCCESS", fmt.Sprintf("Booking %s confirmed", bookingID), booking.UserID)
	s.broadcastSeats(showtimeID, booking.Seats, model.SeatStatusBooked, booking.UserID, bookingID, time.Time{})

	return booking, nil
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID string) ([]model.Booking, error) {
	return s.bookingRepo.FindByUserID(ctx, userID)
}

func (s *bookingService) GetAllBookings(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error) {
	return s.bookingRepo.FindAll(ctx, filter)
}

func (s *bookingService) StartExpirationWorker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bookings, err := s.bookingRepo.FindExpiredPending(context.Background(), time.Now())
			if err != nil {
				s.publishAuditLog(context.Background(), "SYSTEM_ERROR", fmt.Sprintf("Expiration worker failed: %v", err), "")
				continue
			}
			for i := range bookings {
				_ = s.expireBooking(context.Background(), &bookings[i])
			}
		}
	}
}

func (s *bookingService) scheduleBookingExpiration(bookingID string, delay time.Duration) {
	if delay < 0 {
		delay = 0
	}

	time.AfterFunc(delay, func() {
		booking, err := s.bookingRepo.FindByID(context.Background(), bookingID)
		if err != nil || booking == nil || booking.Status != model.BookingStatusPending {
			return
		}
		if time.Now().Before(booking.LockExpiresAt) {
			s.scheduleBookingExpiration(bookingID, time.Until(booking.LockExpiresAt))
			return
		}
		_ = s.expireBooking(context.Background(), booking)
	})
}

func (s *bookingService) expireBooking(ctx context.Context, booking *model.Booking) error {
	if booking == nil || booking.Status != model.BookingStatusPending {
		return nil
	}

	showtimeID := booking.ShowtimeID.Hex()
	for _, seatNo := range booking.Seats {
		_ = s.showtimeRepo.ReleaseSeatLock(ctx, showtimeID, seatNo, booking.UserID)
		_ = s.lockRepo.ReleaseLock(ctx, lockKey(showtimeID, seatNo), booking.LockToken)
	}

	if err := s.bookingRepo.UpdateStatus(ctx, booking.ID.Hex(), model.BookingStatusTimeout); err != nil {
		return err
	}

	s.publishAuditLog(ctx, "BOOKING_TIMEOUT", fmt.Sprintf("Booking %s timed out", booking.ID.Hex()), booking.UserID)
	s.publishAuditLog(ctx, "SEAT_RELEASED", fmt.Sprintf("Released seats %s for booking %s", strings.Join(booking.Seats, ","), booking.ID.Hex()), booking.UserID)
	s.broadcastSeats(showtimeID, booking.Seats, model.SeatStatusAvailable, "", booking.ID.Hex(), time.Time{})
	return nil
}

func (s *bookingService) rollbackLocks(ctx context.Context, showtimeID, userID, lockToken string, acquiredSeats, dbLockedSeats []string) {
	for _, seatNo := range dbLockedSeats {
		_ = s.showtimeRepo.ReleaseSeatLock(ctx, showtimeID, seatNo, userID)
	}
	for _, seatNo := range acquiredSeats {
		_ = s.lockRepo.ReleaseLock(ctx, lockKey(showtimeID, seatNo), lockToken)
	}
	if len(dbLockedSeats) > 0 {
		s.broadcastSeats(showtimeID, dbLockedSeats, model.SeatStatusAvailable, "", "", time.Time{})
	}
}

func (s *bookingService) publishAuditLog(ctx context.Context, eventType, details, userID string) {
	logMsg := model.AuditLog{
		EventType: eventType,
		Details:   details,
		UserID:    userID,
		Timestamp: time.Now(),
	}
	payload, _ := json.Marshal(logMsg)
	_ = s.mqRepo.Publish(ctx, "audit_logs", string(payload))
}

func (s *bookingService) broadcastSeats(showtimeID string, seats []string, status string, lockedBy string, bookingID string, lockedUntil time.Time) {
	if s.wsHub == nil {
		return
	}

	payload := make([]model.Seat, 0, len(seats))
	for _, seatNo := range seats {
		seat := model.Seat{SeatNo: seatNo, Status: status}
		if lockedBy != "" {
			seat.LockedBy = &lockedBy
		}
		if !lockedUntil.IsZero() {
			seat.LockedUntil = &lockedUntil
		}
		payload = append(payload, seat)
	}

	s.wsHub.BroadcastMessage(WSMessage{
		Type:       "SEAT_UPDATED",
		ShowtimeID: showtimeID,
		Payload: map[string]interface{}{
			"booking_id": bookingID,
			"seats":      payload,
		},
	})
}

func validateRequestedSeats(showtime *model.Showtime, requested []string) error {
	bySeat := map[string]model.Seat{}
	for _, seat := range showtime.Seats {
		bySeat[seat.SeatNo] = seat
	}

	now := time.Now()
	for _, seatNo := range requested {
		seat, ok := bySeat[seatNo]
		if !ok {
			return fmt.Errorf("seat %s does not exist", seatNo)
		}
		switch seat.Status {
		case model.SeatStatusBooked:
			return fmt.Errorf("seat %s is already booked", seatNo)
		case model.SeatStatusLocked:
			if seat.LockedUntil != nil && seat.LockedUntil.After(now) {
				return fmt.Errorf("seat %s is currently locked", seatNo)
			}
		}
	}
	return nil
}
