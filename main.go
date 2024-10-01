package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/robfig/cron/v3"
)

type Task struct {
	Name     string
	Func     func(logger slog.Logger)
	Duration string
}

// Task List
var tasks []Task = []Task{
	{Name: "Anemos Cache", Func: writeCache, Duration: "@every 10s"},
}

func main() {
	// Create Logger
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)
	logger := slog.New(
		handler,
	)

	logger.Info("Starting Anemos Public API Collector")

	// Create Scheduler
	logger.Debug("TaskList", slog.Any("tasks", tasks))

	c := cron.New()

	// Add Tasks
	for i := range tasks {
		task := tasks[i]
		ci, err := c.AddFunc(task.Duration, func() { task.Func(*logger.With("funcName", task.Name)) })

		if err != nil {
			logger.Error("Error creating job", "task", task.Name, slog.Any("error", err))
		}

		logger.Info("Task Added", "TaskID:", ci, "Func Name", task.Name, "cron schedule", task.Duration)
	}

	// Start Scheduler
	logger.Info("Start Scheduler")
	logger.Info("If you want to stop the scheduler, please send SIGINT (Ctrl + C)")
	c.Start()

	// シグナル用のチャネル定義
	quit := make(chan os.Signal, 1)

	// 受け取るシグナルを設定
	signal.Notify(quit, os.Interrupt)

	<-quit // ここでシグナルを受け取るまで以降の処理はされない

	// シグナルを受け取った後にしたい処理を書く
	logger.Info("Received SIGINT")
	con := c.Stop()
	logger.Info("Scheduler Stopped", "Context:", con)

}
