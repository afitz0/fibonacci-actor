package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"

	"go.uber.org/zap/zapcore"

	"fibonacci"
	"fibonacci/zapadapter"
)

func main() {
	logger := zapadapter.NewZapAdapter(zapadapter.NewZapLogger(zapcore.DebugLevel))
	c, err := client.Dial(client.Options{
		Logger: logger,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	id := uuid.New().String()
	workflowOptions := client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: "fibonacci",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, fibonacci.Workflow, "Hello", "World")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow asynchronously", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
