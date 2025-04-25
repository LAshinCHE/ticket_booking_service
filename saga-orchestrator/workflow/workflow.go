package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type BookingSagaParams struct {
	UserID   int64
	TicketID int64
	Amount   float64
}

func BookingSagaWorkflow(ctx workflow.Context, params BookingSagaParams) error {

	logger := workflow.GetLogger(ctx)
	logger.Info("Starting BookingSagaWorkflow", "UserID", params.UserID, "TicketID", params.TicketID)

	opts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, opts)

	err = workflow.ExecuteActivity(ctx, ReserveTicketActivity, params).Get(ctx, nil)
	if err != nil {
		_ = workflow.ExecuteActivity(ctx, CancelBookingActivity, params).Get(ctx, nil)
		return err
	}

	err = workflow.ExecuteActivity(ctx, ProcessPaymentActivity, params).Get(ctx, nil)
	if err != nil {
		_ = workflow.ExecuteActivity(ctx, CancelTicketActivity, params).Get(ctx, nil)
		_ = workflow.ExecuteActivity(ctx, CancelBookingActivity, params).Get(ctx, nil)
		return err
	}

	err = workflow.ExecuteActivity(ctx, ConfirmBookingActivity, params).Get(ctx, nil)
	if err != nil {
		return err
	}

	_ = workflow.ExecuteActivity(ctx, SendNotificationActivity, params).Get(ctx, nil)

	return nil
}
