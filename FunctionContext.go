package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
)

type FunctionContext struct {
	Context context.Context
	SpanId  string
	Logger  *zerolog.Logger
}

func FuncCtx() FunctionContext {
	spanId := shortid.MustGenerate()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Info Message")
	log.Error().Msg("Error Message")
	log.Warn().Msg("Warn Message")
	log.Debug().Msg("Debug Message")
	log.Trace().Msg("Trace Message")
	//log.Panic().Msg("Panic Message")
	log.Fatal().Msg("Fatal Message")

	return FunctionContext{
		SpanId: spanId,
	}

}
