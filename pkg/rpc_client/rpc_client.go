package rpc_client

import (
	"log"
	"net/rpc"
	"time"
)

type BleveService struct {
	Client        *rpc.Client
	RetryLimit    uint8
	retryCount    uint8
	DelayDuration time.Duration
	addr          string
}

func New(addr string, retryLimit uint8, delayDuration time.Duration) *BleveService {
	if addr == "" {
		log.Fatal("RPC CLIENT ERROR | BleveService | Address not set for this foundation service!")
	}

	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Println("RPC CLIENT ERROR | BleveService | Dialing TCP Error:", err)
		return nil
	}

	return &BleveService{
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
			log.Println("RPC CLIENT ERROR | BleveService | Detected 'connection is shut down' | Retrying #", s.retryCount)

			// We need to apply an artifical delay in case we need to give time
			// for the server is starting up.
			time.Sleep(s.DelayDuration)

			// Attempt to re-connected and if successful then retry calling the
			// RPC endpoint, else return with error.
			client, err := rpc.DialHTTP("tcp", s.addr)
			if err != nil {
				log.Println("RPC CLIENT ERROR | BleveService | Detected 'connection is shut down' | Failed reconnecting | err:", err.Error())

				// Note: Use recursion to retry the call.
				return s.call(serviceMethod, args, reply)
			}

			log.Println("RPC CLIENT ERROR | BleveService | Detected 'connection is shut down' | Reconnected!")
			s.Client = client

			// Note: Use recursion to retry the call.
			return s.call(serviceMethod, args, reply)
		}
		log.Println("RPC CLIENT ERROR | BleveService | Detected 'connection is shut down' | Too many retries | err:", err.Error())
		return err
	}

	// If success then nil will be returned, else the specific error will be returned.
	return err
}
