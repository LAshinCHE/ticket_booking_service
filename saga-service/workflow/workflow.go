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
		"BookingID", input.BookingData.ID,
		"UserID", input.BookingData.UserID,
		"TicketID", input.BookingData.TicketID,
		"Price", input.BookingData.Price,
		"TraceCtx", input.TraceCtx,
	)

	// Настраиваем Activity-опции (таймауты / ретраи при желании)
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	})
	var err error
	// Нужен указатель, чтобы Temporal нашёл методы
	var acts *activities.BookingActivities
	var createBookingResult activities.CreateBookingResulet
	var updatedTraceCtx map[string]string

	// Компенсация на случай дальнейших сбоев
	defer func() {
		if err != nil {
			workflow.ExecuteActivity(ctx,
				acts.CancelBooking,
				input.BookingData.ID,
				updatedTraceCtx,
			).Get(ctx, nil)
		}
	}()

	//----------------------------------------------------------------------
	// 1. Создаём бронирование
	if err = workflow.ExecuteActivity(ctx,
		acts.CreateBooking, input.BookingData.ID, input.BookingData.UserID, input.BookingData.TicketID, input.BookingData.Price, input.TraceCtx).
		Get(ctx, &createBookingResult); err != nil {
		return err
	}
	updatedTraceCtx = createBookingResult.TraceCtx

	defer func() {
		if err != nil {
			workflow.ExecuteActivity(ctx,
				acts.MakeAvailableTicket,
				input.BookingData.TicketID,
				updatedTraceCtx).Get(ctx, nil)
		}
	}()

	//----------------------------------------------------------------------
	// 2. Резервируем билет
	if err = workflow.ExecuteActivity(ctx,
		acts.ReserveTicket,
		input.BookingData.TicketID,
		updatedTraceCtx,
	).Get(ctx, &updatedTraceCtx); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			workflow.ExecuteActivity(ctx,
				acts.CancelWithdrawMoney,
				input.BookingData.UserID,
				input.BookingData.Price,
				updatedTraceCtx).Get(ctx, nil)
		}
	}()
	//----------------------------------------------------------------------
	// 4. Списываем деньги
	if err = workflow.ExecuteActivity(ctx,
		acts.WithdrawMoney,
		input.BookingData.UserID,
		input.BookingData.Price,
		updatedTraceCtx).Get(ctx, nil); err != nil {
		return err
	}

	// 5. Уведомляем пользователя (не критично, поэтому без проверки Get)
	err = workflow.ExecuteActivity(ctx,
		acts.NotifyUser,
		input.BookingData.UserID,
		"Бронирование успешно",
		updatedTraceCtx).Get(ctx, nil)

	if err != nil {
		logger.Warn("Failed to notify user", "err", err)
	}

	logger.Info("Saga finished OK", "BookingID", createBookingResult.ID)
	return nil
}
