package rpc_client

import (
	"os"
	"net/rpc"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog"
)

type BleveService struct {
	Logger        zerolog.Logger
	Client        *rpc.Client
	RetryLimit    uint8
	retryCount    uint8
	DelayDuration time.Duration
	addr          string
}

func New(addr string, retryLimit uint8, delayDuration time.Duration) *BleveService {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if addr == "" {
		logger.Fatal().Caller().Msg("address not set for this bleve service")
	}

	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		logger.Fatal().Err(err).Caller().Msg("dialing tcp error")
		return nil
	}

	return &BleveService{
		Logger:        logger,
		Client:        client,
		RetryLimit:    retryLimit,
		retryCount:    0,
		DelayDuration: delayDuration,
		addr:          addr,
	}
}

// Function used to make RPC calls with retry functionality in case the RPC
// server has been shutdown and the connection was lost.
func (s *BleveService) call(serviceMethod string, args interface{}, reply interface{}) error {
	err := s.Client.Call(serviceMethod, args, reply)

	// Detect the `connection is shut down` error.
	if err == rpc.ErrShutdown {
		if s.retryCount < s.RetryLimit {
			s.retryCount += 1
            s.Logger.Warn().Caller().Msgf("detected 'connection is shut down' so retrying #%d", s.retryCount)

			// We need to apply an artifical delay in case we need to give time
			// for the server is starting up.
			time.Sleep(s.DelayDuration)

			// Attempt to re-connected and if successful then retry calling the
			// RPC endpoint, else return with error.
			client, err := rpc.DialHTTP("tcp", s.addr)
			if err != nil {
				s.Logger.Error().Err(err).Caller().Msg("detected 'connection is shut down' and failed reconnecting")

				// Note: Use recursion to retry the call.
				return s.call(serviceMethod, args, reply)
			}

			s.Logger.Warn().Caller().Msg("detected 'connection is shut down' so reconnected")
			s.Client = client

			// Note: Use recursion to retry the call.
			return s.call(serviceMethod, args, reply)
		}

		s.Logger.Error().Err(err).Caller().Msg("detected 'connection is shut down' and had too many retries")
		return err
	}

	// If success then nil will be returned, else the specific error will be returned.
	return err
}
