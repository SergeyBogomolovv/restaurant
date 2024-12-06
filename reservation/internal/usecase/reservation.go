package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/SergeyBogomolovv/restaurant/common/constants"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/errors"
	"github.com/google/uuid"
)

type ReservationRepo interface {
	CreateReservation(ctx context.Context, dto *dto.CreateReservationDTO) (*dto.ReservationCreated, error)
	SetReservationStatus(ctx context.Context, reservationID uuid.UUID, status string) error
	CloseEndedReservations(ctx context.Context) (int64, error)
	GetTableExists(ctx context.Context, tableID uuid.UUID) (bool, error)
}

type ReservationBroker interface {
	Publish(key string, payload any) error
}

type reservationUsecase struct {
	log    *slog.Logger
	repo   ReservationRepo
	broker ReservationBroker
}

func NewReservationUsecase(log *slog.Logger, repo ReservationRepo, broker ReservationBroker) *reservationUsecase {
	usecase := &reservationUsecase{log: log, repo: repo, broker: broker}
	return usecase
}

func (u *reservationUsecase) RunEndedReservationsChecker(ctx context.Context, duration time.Duration) {
	const op = "reservation.CheckEndedReservations"
	log := u.log.With(slog.String("op", op))

	log.Info("reservations checker started")

	for {
		now := time.Now()
		next := now.Truncate(duration).Add(duration)
		waitDuration := time.Until(next)

		timer := time.NewTimer(waitDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Info("reservations checker stopped")
			return
		case <-timer.C:
			count, err := u.repo.CloseEndedReservations(ctx)
			if err != nil {
				log.Error("failed to check closed reservations", "error", err)
			}
			if count > 0 {
				log.Info("closed reservations", "count", count)
			}
		}
	}
}

func (u *reservationUsecase) CreateReservation(ctx context.Context, payload *dto.CreateReservationDTO) (uuid.UUID, error) {
	const op = "reservation.Create"
	log := u.log.With(slog.String("op", op))

	log.Info("creating reservation")

	tableExists, err := u.repo.GetTableExists(ctx, payload.TableID)
	if err != nil {
		log.Error("failed to check table exists", "error", err)
		return uuid.Nil, err
	}
	if !tableExists {
		log.Info("table not found")
		return uuid.Nil, errs.ErrTableNotFound
	}

	reservation, err := u.repo.CreateReservation(ctx, payload)
	if err != nil {
		if errors.Is(err, errs.ErrTableAlreadyReserved) {
			log.Info("table already reserved")
			return uuid.Nil, errs.ErrTableAlreadyReserved
		}
		log.Error("failed to create reservation", "error", err)
		return uuid.Nil, err
	}

	if err := u.broker.Publish("reservation.created", reservation); err != nil {
		log.Error("failed to publish message", "error", err)
		return uuid.Nil, err
	}

	return reservation.ReservationID, nil
}

func (u *reservationUsecase) CancelReservation(ctx context.Context, reservationId uuid.UUID) error {
	const op = "reservation.Cancel"
	log := u.log.With(slog.String("op", op))

	log.Info("cancelling reservation")

	//TODO: check is admin or current user

	if err := u.repo.SetReservationStatus(ctx, reservationId, constants.ReservationStatusCancelled); err != nil {
		if errors.Is(err, errs.ErrReservationNotFound) {
			log.Info("reservation not found")
			return errs.ErrReservationNotFound
		}
		log.Error("failed to cancel reservation", "error", err)
		return err
	}

	//TODO: publish message

	return nil
}

func (u *reservationUsecase) CloseReservation(ctx context.Context, reservationId uuid.UUID) error {
	const op = "reservation.Close"
	log := u.log.With(slog.String("op", op))

	log.Info("closing reservation")

	//TODO: check is admin or waiter

	if err := u.repo.SetReservationStatus(ctx, reservationId, constants.ReservationStatusClosed); err != nil {
		if errors.Is(err, errs.ErrReservationNotFound) {
			log.Info("reservation not found")
			return errs.ErrReservationNotFound
		}
		log.Error("failed to close reservation", "error", err)
		return err
	}

	//TODO: publish message

	return nil
}
