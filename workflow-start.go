package fibonacci

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// 52 is the sweet spot on my 2021 M1 Max macbook, clocking in at just under 2 minutes
// (assuming no other CPU contention at the time).
const MaxFibonacciNumber = 52

func Workflow(ctx workflow.Context, greeting string, name string) (err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 120 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Fibonacci workflow started", "greeting", greeting, "name", name)

	selector := workflow.NewSelector(ctx)

	fibNumber := 0
	for fibNumber < MaxFibonacciNumber {
		fibFuture := workflow.ExecuteActivity(ctx, FibonacciActivity, fibNumber)
		selector.AddFuture(fibFuture, func(f workflow.Future) {
			var result int
			actErr := f.Get(ctx, &result)
			if actErr != nil {
				logger.Error("Activity failed.", "Error", err)
				err = actErr
				return
			}
			logger.Info("Finished calculating fibonacci number", "Number", fibNumber, "Result", result)
			fibNumber++
		})

		selector.Select(ctx)
		if err != nil {
			return err
		}
	}

	logger.Info("Workflow completed.")
	return nil
}
