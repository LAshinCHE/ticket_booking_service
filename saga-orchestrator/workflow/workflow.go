package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/activities"
)

// BookingSagaWorkflow — сам воркфлоу.
func BookingSagaWorkflow(ctx workflow.Context, p BookingParams) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Saga started",
		"UserID", p.UserID, "TicketID", p.TicketID, "Price", p.Price)

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
		acts.CreateBooking, p.UserID, p.TicketID, p.Price).
		Get(ctx, &bookingID); err != nil {
		return err
	}
	// дефолтная компенсация на случай любого дальнейшего сбоя
	defer workflow.ExecuteActivity(ctx, acts.CancelBooking, bookingID).Get(ctx, nil)

	//----------------------------------------------------------------------
	// 2. Проверяем наличие билета
	var ok bool
	if err := workflow.ExecuteActivity(ctx,
		acts.CheckTicketAvailability, p.TicketID).Get(ctx, &ok); err != nil {
		return err
	}
	if !ok {
		return workflow.NewContinueAsNewError(ctx, BookingSagaWorkflow, p) // или любая ваша ошибка
	}

	//----------------------------------------------------------------------
	// 3. Резервируем билет
	if err := workflow.ExecuteActivity(ctx,
		acts.ReserveTicket, p.TicketID).Get(ctx, nil); err != nil {
		return err
	}

	//----------------------------------------------------------------------
	// 4. Списываем деньги
	if err := workflow.ExecuteActivity(ctx,
		acts.WithdrawMoney, p.UserID, p.Price, bookingID).Get(ctx, nil); err != nil {
		return err
	}
	defer workflow.ExecuteActivity(ctx, acts.CancelWithdrawMoney, p.UserID, p.Price).Get(ctx, nil)
	//----------------------------------------------------------------------
	// 5. Уведомляем пользователя (не критично, поэтому без проверки Get)
	_ = workflow.ExecuteActivity(ctx,
		acts.NotifyUser, p.UserID, "Бронирование успешно").Get(ctx, nil)

	logger.Info("Saga finished OK", "BookingID", bookingID)
	return nil
}
