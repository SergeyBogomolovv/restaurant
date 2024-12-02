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
		return nil, status.Error(codes.InvalidArgument, "invalid customerId")
	}
	tableId, err := uuid.Parse(req.TableId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid tableId")
	}

	dto := &dto.CreateReservationDTO{
		CustomerID: customerId,
		TableID:    tableId,
		StartTime:  time.Unix(req.StartTime, 0),
		EndTime:    time.Unix(req.EndTime, 0),
	}
	if dto.StartTime.After(dto.EndTime) || time.Now().After(dto.StartTime) {
		return nil, status.Error(codes.InvalidArgument, "invalid time range")
	}

	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload: %v", err)
	}

	reservationId, err := h.reservationUsecase.CreateReservation(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrTableAlreadyReserved):
			return nil, status.Error(codes.AlreadyExists, "table already reserved")
		case errors.Is(err, errs.ErrTableNotFound):
			return nil, status.Error(codes.NotFound, "table not found")
		default:
			return nil, status.Error(codes.Internal, "failed to create reservation")
		}
	}

	return &pb.CreateReservationResponse{ReservationId: reservationId.String()}, nil
}

func (h *reservationHandler) CancelReservation(ctx context.Context, req *pb.CancelReservationRequest) (*pb.CancelReservationResponse, error) {
	reservationId, err := uuid.Parse(req.ReservationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reservationId")
	}

	if err := h.reservationUsecase.CancelReservation(ctx, reservationId); err != nil {
		switch {
		case errors.Is(err, errs.ErrReservationNotFound):
			return nil, status.Error(codes.NotFound, "reservation not found")
		default:
			return nil, status.Error(codes.Internal, "failed to cancel reservation")
		}
	}

	return &pb.CancelReservationResponse{Status: "cancelled"}, nil
}

func (h *reservationHandler) CloseReservation(ctx context.Context, req *pb.CloseReservationRequest) (*pb.CloseReservationResponse, error) {
	reservationId, err := uuid.Parse(req.ReservationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reservationId")
	}

	if err := h.reservationUsecase.CloseReservation(ctx, reservationId); err != nil {
		switch {
		case errors.Is(err, errs.ErrReservationNotFound):
			return nil, status.Error(codes.NotFound, "reservation not found")
		default:
			return nil, status.Error(codes.Internal, "failed to close reservation")
		}
	}

	return &pb.CloseReservationResponse{Status: "closed"}, nil
}
