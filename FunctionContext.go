package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/teris-io/shortid"
	"os"
)

type FunctionContext struct {
	Context        context.Context
	SpanId         string
	spanIdLogField string
	Logger         *zerolog.Logger
}

func FuncCtx(ctx context.Context) FunctionContext {
	spanId := shortid.MustGenerate()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Ctx(ctx).Timestamp().Str("spanId", "["+spanId+"]").Logger()

	logger.Info().Msg("Info Message")
	logger.Error().Msg("Error Message")
	logger.Warn().Msg("Warn Message")
	logger.Debug().Msg("Debug Message")
	logger.Trace().Msg("Trace Message")
	//log.Panic().Msg("Panic Message")
	//logger.Fatal().Msg("Fatal Message")

	return FunctionContext{
		SpanId:         spanId,
		spanIdLogField: "[" + spanId + "] ",
		Logger:         &logger,
	}

}
