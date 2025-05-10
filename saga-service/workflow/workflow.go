package workflow

import (
	"context"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/LAshinCHE/ticket_booking_service/saga-service/activities"
	"github.com/LAshinCHE/ticket_booking_service/saga-service/metrics"
)

func BookingSagaWorkflow(ctx workflow.Context, input BookingWorkflowInput) error {
	workflow.SideEffect(ctx, func(workflow.Context) interface{} {
		metrics.IncSagaStarted(context.Background())
		return nil
	})
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

	defer func() {
		if err != nil {
			workflow.ExecuteActivity(ctx,
				acts.CancelBooking,
				input.BookingData.ID,
				updatedTraceCtx,
			).Get(ctx, nil)
			workflow.SideEffect(ctx, func(workflow.Context) interface{} {
				metrics.IncSagaFailed(context.Background(), err)
				return nil
			})
		}
	}()
	defer func() {
		if err == nil {
			_ = workflow.ExecuteActivity(ctx,
				acts.NotifyUser,
				input.BookingData.UserID,
				"Бронирование успешно",
				updatedTraceCtx).Get(ctx, nil)
		} else {
			_ = workflow.ExecuteActivity(ctx,
				acts.NotifyUser,
				input.BookingData.UserID,
				"Бронирование завершилось с ошибкой "+err.Error(),
				updatedTraceCtx).Get(ctx, nil)
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
	workflow.SideEffect(ctx, func(workflow.Context) interface{} {
		metrics.IncSagaSucceeded(context.Background())
		return nil
	})
	logger.Info("Saga finished OK", "BookingID", createBookingResult.ID)
	return nil
}
