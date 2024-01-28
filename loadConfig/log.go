package loadConfig

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"log"
	"log/slog"
	"os"
)

func LogInit() {

	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	file, err := os.OpenFile("log/Server.log", os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	handler := NewPrettyHandler(os.Stdout, file, opts)
	//handler := NewPrettyHandler(os.Stdout, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	//slog.Debug(
	//	"executing database query",
	//	slog.String("query", "SELECT * FROM users"),
	//)
	//slog.Info("image upload successful", slog.String("image_id", "39ud88"))
	//slog.Warn("Error Network", slog.String("s", "ds"))
}

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	//l *slog.Logger
	stdOutLog *log.Logger
	fileLog   *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	fileLevel := level

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[2006/01/02 15:05:05.000]")
	msg := color.CyanString(r.Message)

	// 分别写文件和打印到控制台
	h.stdOutLog.Println(timeStr, level, msg, color.WhiteString(string(b)))
	h.fileLog.Println(timeStr, fileLevel, r.Message, string(b))

	return nil
}
func NewPrettyHandler(out io.Writer, file io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler:   slog.NewJSONHandler(out, &opts.SlogOpts),
		stdOutLog: log.New(out, "", 0),
		fileLog:   log.New(file, "", 0),
	}

	return h
}
