package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"go.uber.org/zap/zapcore"

	"fibonacci"
	"fibonacci/zapadapter"
)

func main() {
	logger := zapadapter.NewZapAdapter(zapadapter.NewZapLogger(zapcore.InfoLevel))
	c, err := client.Dial(client.Options{
		Logger: logger,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "fibonacci", worker.Options{})

	w.RegisterWorkflow(fibonacci.Workflow)
	w.RegisterWorkflow(fibonacci.WorkflowFinal)
	w.RegisterWorkflow(fibonacci.WorkflowStart)
	w.RegisterActivity(fibonacci.FibonacciActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
