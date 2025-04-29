package workflow

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.temporal.io/sdk/workflow"

	"github.com/LAshinCHE/ticket_booking_service/saga-service/activities"
)

func BookingSagaWorkflow(ctx workflow.Context, input BookingWorkflowInput) error {

	propagator := propagation.TraceContext{}
	traceCarrier := propagation.MapCarrier(input.TraceCtx)
	realCtx := propagator.Extract(context.Background(), traceCarrier)

	_, span := otel.Tracer("saga-service").Start(realCtx, "BookingSagaWorkflow")
	defer span.End()

	logger := workflow.GetLogger(ctx)
	logger.Info("Saga started",
		"UserID", input.Params.UserID, "TicketID", input.Params.TicketID, "Price", input.Params.Price)

	// Настраиваем Activity-опции (таймауты / ретраи при желании)
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	// Нужен указатель, чтобы Temporal нашёл методы
	var acts *activities.BookingActivities

	//----------------------------------------------------------------------
	// 1. Создаём бронирование
	var bookingID string
	if err := workflow.ExecuteActivity(ctx,
		acts.CreateBooking, input.Params.UserID, input.Params.TicketID, input.Params.Price).
		Get(ctx, &bookingID); err != nil {
		return err
	}
	// дефолтная компенсация на случай любого дальнейшего сбоя
	defer workflow.ExecuteActivity(ctx, acts.CancelBooking, bookingID).Get(ctx, nil)

	//----------------------------------------------------------------------
	// 2. Проверяем наличие билета
	var ok bool
	if err := workflow.ExecuteActivity(ctx,
		acts.CheckTicketAvailability, input.Params.TicketID).Get(ctx, &ok); err != nil {
		return err
	}
	if !ok {
		return workflow.NewContinueAsNewError(ctx, BookingSagaWorkflow, input) // или любая ваша ошибка
	}

	//----------------------------------------------------------------------
	// 3. Резервируем билет
	if err := workflow.ExecuteActivity(ctx,
		acts.ReserveTicket, input.Params.TicketID).Get(ctx, nil); err != nil {
		return err
	}

	//----------------------------------------------------------------------
	// 4. Списываем деньги
	if err := workflow.ExecuteActivity(ctx,
		acts.WithdrawMoney, input.Params.UserID, input.Params.Price, bookingID).Get(ctx, nil); err != nil {
		return err
	}
	defer workflow.ExecuteActivity(ctx, acts.CancelWithdrawMoney, input.Params.UserID, input.Params.Price).Get(ctx, nil) // возвращаем деньги в случае ошибки

	// 5. Уведомляем пользователя (не критично, поэтому без проверки Get)
	_ = workflow.ExecuteActivity(ctx,
		acts.NotifyUser, input.Params.UserID, "Бронирование успешно").Get(ctx, nil)

	logger.Info("Saga finished OK", "BookingID", bookingID)
	return nil
}
