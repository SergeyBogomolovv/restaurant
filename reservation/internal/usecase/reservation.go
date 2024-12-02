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

type Repo interface {
	CreateReservation(ctx context.Context, dto *dto.CreateReservationDTO) (uuid.UUID, error)
	SetReservationStatus(ctx context.Context, reservationID uuid.UUID, status string) error
	CloseEndedReservations(ctx context.Context) (int64, error)
	GetTableExists(ctx context.Context, tableID uuid.UUID) (bool, error)
}

type reservationUsecase struct {
	log  *slog.Logger
	repo Repo
}

func NewReservationUsecase(log *slog.Logger, repo Repo, ctx context.Context, tickerDuration time.Duration) *reservationUsecase {
	usecase := &reservationUsecase{log: log, repo: repo}
	go usecase.CheckEndedReservations(ctx, tickerDuration)
	return usecase
}

func (u *reservationUsecase) CreateReservation(ctx context.Context, dto *dto.CreateReservationDTO) (uuid.UUID, error) {
	const op = "reservation.Create"
	log := u.log.With(slog.String("op", op))

	log.Info("creating reservation")

	tableExists, err := u.repo.GetTableExists(ctx, dto.TableID)
	if err != nil {
		log.Error("failed to check table exists", "error", err)
		return uuid.Nil, err
	}
	if !tableExists {
		log.Info("table not found")
		return uuid.Nil, errs.ErrTableNotFound
	}

	id, err := u.repo.CreateReservation(ctx, dto)
	if err != nil {
		if errors.Is(err, errs.ErrTableAlreadyReserved) {
			log.Info("table already reserved")
			return uuid.Nil, errs.ErrTableAlreadyReserved
		}
		log.Error("failed to create reservation", "error", err)
		return uuid.Nil, err
	}

	return id, nil
}

func (u *reservationUsecase) CheckEndedReservations(ctx context.Context, duration time.Duration) {
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
			log.Info("check closed reservations stopped")
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

	// TODO: send rmq

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

	// TODO: send rmq

	return nil
}
