package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/LAshinCHE/ticket_booking_service/saga-service/activities"
)

func BookingSagaWorkflow(ctx workflow.Context, input BookingWorkflowInput) error {

	logger := workflow.GetLogger(ctx)
	logger.Info("Saga started",
		"BookingID", input.BookingData.ID, "UserID", input.BookingData.UserID, "TicketID", input.BookingData.TicketID, "Price", input.BookingData.Price, "TraceTx", input.TraceCtx)

	// Настраиваем Activity-опции (таймауты / ретраи при желании)
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	})

	// Нужен указатель, чтобы Temporal нашёл методы
	var acts *activities.BookingActivities

	//----------------------------------------------------------------------
	// 1. Создаём бронирование
	var bookingID string
	if err := workflow.ExecuteActivity(ctx,
		acts.CreateBooking, input.BookingData.ID, input.BookingData.UserID, input.BookingData.TicketID, input.BookingData.Price, input.TraceCtx).
		Get(ctx, &bookingID); err != nil {
		return err
	}
	// дефолтная компенсация на случай любого дальнейшего сбоя
	defer workflow.ExecuteActivity(ctx, acts.CancelBooking, bookingID).Get(ctx, nil)

	//----------------------------------------------------------------------
	// 2. Проверяем наличие билета
	var ok bool
	if err := workflow.ExecuteActivity(ctx,
		acts.CheckTicketAvailability, input.BookingData.TicketID).Get(ctx, &ok); err != nil {
		return err
	}
	if !ok {
		return workflow.NewContinueAsNewError(ctx, BookingSagaWorkflow, input) // или любая ваша ошибка
	}

	//----------------------------------------------------------------------
	// 3. Резервируем билет
	if err := workflow.ExecuteActivity(ctx,
		acts.ReserveTicket, input.BookingData.TicketID).Get(ctx, nil); err != nil {
		return err
	}

	//----------------------------------------------------------------------
	// 4. Списываем деньги
	if err := workflow.ExecuteActivity(ctx,
		acts.WithdrawMoney, input.BookingData.UserID, input.BookingData.Price, bookingID).Get(ctx, nil); err != nil {
		return err
	}
	defer workflow.ExecuteActivity(ctx, acts.CancelWithdrawMoney, input.BookingData.UserID, input.BookingData.Price).Get(ctx, nil) // возвращаем деньги в случае ошибки

	// 5. Уведомляем пользователя (не критично, поэтому без проверки Get)
	_ = workflow.ExecuteActivity(ctx,
		acts.NotifyUser, input.BookingData.UserID, "Бронирование успешно").Get(ctx, nil)

	logger.Info("Saga finished OK", "BookingID", bookingID)
	return nil
}
