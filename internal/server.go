package tracker

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/uber/jaeger-client-go/config"
)

type tracer struct {
	Tracer opentracing.Tracer
	Closer io.Closer
}

type Server struct {
	tasks []Task
	zerolog.Logger
	CounterVec   *prometheus.CounterVec
	HistogramVec *prometheus.CounterVec
	tracer       *tracer
}

func initTracer(serviceName string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			CollectorEndpoint: "http://jaeger:14268/api/traces",
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic(err)
	}

	return tracer, closer
}

func initPrometheus(requestsTotal, requestDuration *prometheus.CounterVec) {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
}

func NewServer(logFile *os.File, requestsTotal, requestDuration *prometheus.CounterVec) *Server {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(logFile).With().Timestamp().Logger()

	initPrometheus(requestsTotal, requestDuration)

	initTracer, closer := initTracer("tracker")


	return &Server{
		tasks:        make([]Task, 0),
		Logger:       logger,
		CounterVec:   requestsTotal,
		HistogramVec: requestDuration,
		tracer:       &tracer{
			Tracer: initTracer,
			Closer: closer,
		},
	}
}

func randSleep() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

// Метод создает задание с заданным id в образовательной программе
// (DELETE /task/{idTask})
func (s *Server) DeleteTaskIdTask(ctx echo.Context, idTask int) error {
	rootSpan := s.tracer.Tracer.StartSpan("DeleteTaskIdTask")
	defer rootSpan.Finish()

	rootSpan.SetTag("method", "delete")

	timeStart := time.Now()

	for i, task := range s.tasks {
		span1 := s.tracer.Tracer.StartSpan("FindTaskIdTask", opentracing.ChildOf(rootSpan.Context()))
		defer span1.Finish()
		if task.IdTask == idTask {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			
			span2 := s.tracer.Tracer.StartSpan("Logging about delete task", opentracing.ChildOf(span1.Context()))
			s.Logger.Info().Msg("Task deleted")
			span2.Finish()

			return ctx.JSON(http.StatusOK, task)
		}
	}

	timeEnd := time.Now()
	result := timeEnd.Sub(timeStart)

	span3 := s.tracer.Tracer.StartSpan("Sleep", opentracing.ChildOf(rootSpan.Context()))
	randSleep()
	span3.Finish()
	s.CounterVec.WithLabelValues("delete").Inc()
	s.HistogramVec.WithLabelValues("delete").Add(result.Seconds())

	return ctx.JSON(http.StatusNotFound, nil)
}

// Метод получает задание по id образовательной программы
// (GET /task/{idTask})
func (s *Server) GetTaskIdTask(ctx echo.Context, idTask int) error {
	rootSpan := s.tracer.Tracer.StartSpan("GetTaskIdTask")
	defer rootSpan.Finish()
	rootSpan.SetTag("method", "get")

	sleepSpan := s.tracer.Tracer.StartSpan("Sleep", opentracing.ChildOf(rootSpan.Context()))
	randSleep()
	sleepSpan.Finish()

	timeStart := time.Now()

	for _, task := range s.tasks {
		span1 := s.tracer.Tracer.StartSpan("FindTaskForIdTask", opentracing.ChildOf(rootSpan.Context()))
		defer span1.Finish()
		if task.IdTask == idTask {
			s.Logger.Info().Msg("Task found")

			return ctx.JSON(http.StatusOK, task)
		}
	}

	timeEnd := time.Now()
	result := timeEnd.Sub(timeStart)

	s.CounterVec.WithLabelValues("get").Inc()
	s.HistogramVec.WithLabelValues("get").Add(result.Seconds())

	return ctx.JSON(http.StatusNotFound, nil)
}

// Метод получает все задания образовательной программы
// (GET /tasks)
func (s *Server) GetTasks(ctx echo.Context) error {
	rootSpan := s.tracer.Tracer.StartSpan("GetTasks")
	defer rootSpan.Finish()

	sleepSpan := s.tracer.Tracer.StartSpan("Sleep", opentracing.ChildOf(rootSpan.Context()))
	randSleep()
	sleepSpan.Finish()

	timeStart := time.Now()
	timeEnd := time.Now()
	result := timeEnd.Sub(timeStart)

	s.Logger.Info().Msg("Tasks found")
	s.HistogramVec.WithLabelValues("get").Add(result.Seconds())
	s.CounterVec.WithLabelValues("get").Inc()

	return ctx.JSON(http.StatusOK, s.tasks)
}

// Метод создает задание с заданным id в образовательной программе
// (POST /tasks)
func (s *Server) PostTasks(ctx echo.Context) error {
	rootSpan := s.tracer.Tracer.StartSpan("PostTasks")
	defer rootSpan.Finish()

	rootSpan.SetTag("method", "post")

	sleepSpan := s.tracer.Tracer.StartSpan("Sleep", opentracing.ChildOf(rootSpan.Context()))
	randSleep()
	sleepSpan.Finish()

	timeStart := time.Now()

	readSpan := s.tracer.Tracer.StartSpan("ReadBody", opentracing.ChildOf(rootSpan.Context()))
	bytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	readSpan.Finish()

	task := &Task{}
	err = json.Unmarshal(bytes, &task)
	if err != nil {
		return err
	}

	s.tasks = append(s.tasks, *task)

	s.Logger.Info().Msg("Task created")

	timeEnd := time.Now()
	result := timeEnd.Sub(timeStart)

	s.CounterVec.WithLabelValues("post").Inc()
	s.HistogramVec.WithLabelValues("post").Add(float64(result.Nanoseconds()))

	return ctx.JSON(http.StatusOK, task)
}
