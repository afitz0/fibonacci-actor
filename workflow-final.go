package fibonacci

import (
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

// 52 is the sweet spot on my 2021 M1 Max macbook, clocking in at around 2 minutes.
const MaxFibonacciNumber_final = 52

func WorkflowFinal(ctx workflow.Context, startNum int, children []string) (err error) {
	if startNum >= MaxFibonacciNumber_final {
		return nil
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 120 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Fibonacci workflow started", "StartingNumber", startNum)

	selector := workflow.NewSelector(ctx)
	signalCreateNewChannel := workflow.GetSignalChannel(ctx, "signalCreateNew")

	var childIds []string
	selector.AddReceive(signalCreateNewChannel, func(c workflow.ReceiveChannel, _ bool) {
		c.Receive(ctx, nil)
		logger.Info("Signal received")

		childId, childErr := startChild(ctx)
		if childErr != nil {
			logger.Error("Failed to start child workflow")
			err = childErr
			return
		}
		childIds = append(childIds, childId)
	})

	signalHelloChannel := workflow.GetSignalChannel(ctx, "signalHello")
	selector.AddReceive(signalHelloChannel, func(c workflow.ReceiveChannel, _ bool) {
		var senderId string
		c.Receive(ctx, &senderId)
		logger.Info("Hello signal received", "SenderId", senderId)
	})

	fibNumber := startNum
	done := false
	for fibNumber < MaxFibonacciNumber && !done {
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

		// Continue-As-New every 10 numbers
		if fibNumber - startNum > 10 {
			done = true
		}
	}

	// Drain hello signal channel asynchronously to avoid signal loss
	for {
		var senderId string
		ok := signalHelloChannel.ReceiveAsync(&senderId)
		if !ok {
			break
		}
		logger.Info("Hello signal received", "SenderId", senderId)
	}

	// Drain spawn new signal channel asynchronously to avoid signal loss
	for {
		ok := signalHelloChannel.ReceiveAsync(nil)
		if !ok {
			break
		}
		logger.Info("Start Child signal received")

		childId, childErr := startChild(ctx)
		if childErr != nil {
			logger.Error("Failed to start child workflow")
			err = childErr
			return
		}
		childIds = append(childIds, childId)
	}
	return workflow.NewContinueAsNewError(ctx, WorkflowFinal, fibNumber, childIds)
}

func startChild(ctx workflow.Context) (childId string, err error) {
	childId = uuid.New().String()
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		WorkflowID:        childId,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	err = workflow.ExecuteChildWorkflow(ctx, WorkflowFinal, 0, []string{}).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Failed to start child")
		return "", err
	}
	workflow.GetLogger(ctx).Info("Started child workflow", "ChildWorkflowId", childId)

	return childId, nil
}
