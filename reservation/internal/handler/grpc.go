package handler

import (
	"context"
	"errors"
	"time"

	pb "github.com/SergeyBogomolovv/restaurant/common/api/gen/reservation"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReservationUsecase interface {
	CreateReservation(ctx context.Context, dto *dto.CreateReservationDTO) (uuid.UUID, error)
	CancelReservation(ctx context.Context, reservationId uuid.UUID) error
	CloseReservation(ctx context.Context, reservationId uuid.UUID) error
}

type reservationHandler struct {
	validate           *validator.Validate
	reservationUsecase ReservationUsecase
	pb.UnimplementedReservationServer
}

func RegisterGRPCHandler(server *grpc.Server, reservationUsecase ReservationUsecase) {
	handler := &reservationHandler{
		validate:           validator.New(validator.WithRequiredStructEnabled()),
		reservationUsecase: reservationUsecase,
	}
	pb.RegisterReservationServer(server, handler)
}

func (h *reservationHandler) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.CreateReservationResponse, error) {
	customerId, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	tableId, err := uuid.Parse(req.TableId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}

	dto := &dto.CreateReservationDTO{
		CustomerID: customerId,
		TableID:    tableId,
		StartTime:  time.Unix(req.StartTime, 0),
		EndTime:    time.Unix(req.EndTime, 0),
	}
	if dto.StartTime.After(dto.EndTime) || time.Now().After(dto.StartTime) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time range")
	}

	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}

	reservationId, err := h.reservationUsecase.CreateReservation(ctx, dto)
	if err != nil {
		if errors.Is(err, errs.ErrTableAlreadyReserved) {
			return nil, status.Errorf(codes.AlreadyExists, "table already reserved, error: %v", err)
		}
		if errors.Is(err, errs.ErrTableNotFound) {
			return nil, status.Errorf(codes.NotFound, "table not found, error")
		}
		return nil, status.Errorf(codes.Internal, "failed to create reservation, error: %v", err)
	}

	return &pb.CreateReservationResponse{ReservationId: reservationId.String()}, nil
}

func (h *reservationHandler) CancelReservation(ctx context.Context, req *pb.CancelReservationRequest) (*pb.CancelReservationResponse, error) {
	reservationId, err := uuid.Parse(req.ReservationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}

	if err := h.reservationUsecase.CancelReservation(ctx, reservationId); err != nil {
		if errors.Is(err, errs.ErrReservationNotFound) {
			return nil, status.Errorf(codes.NotFound, "reservation not found, error: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel reservation, error: %v", err)
	}

	return &pb.CancelReservationResponse{Status: "cancelled"}, nil
}

func (h *reservationHandler) CloseReservation(ctx context.Context, req *pb.CloseReservationRequest) (*pb.CloseReservationResponse, error) {
	reservationId, err := uuid.Parse(req.ReservationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}

	if err := h.reservationUsecase.CloseReservation(ctx, reservationId); err != nil {
		if errors.Is(err, errs.ErrReservationNotFound) {
			return nil, status.Errorf(codes.NotFound, "reservation not found, error: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to close reservation, error: %v", err)
	}

	return &pb.CloseReservationResponse{Status: "closed"}, nil
}
