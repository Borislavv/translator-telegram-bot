package service

import (
	"context"

	"github.com/Borislavv/Translator-telegram-bot/pkg/api/grpc/service/timeZoneFetcherGRPCInterface"
)

type TimeZoneFetcherGRPC struct {
	userService *UserService
}

// NewTimeZoneFetcherGRPC - constructor of TimeZoneFetcherGRPC structure.
func NewTimeZoneFetcherGRPC(userService *UserService) *TimeZoneFetcherGRPC {
	return &TimeZoneFetcherGRPC{
		userService: userService,
	}
}

// GetTimeZone - proxy method to UserService for determine timezone
func (fetcher *TimeZoneFetcherGRPC) GetTimeZone(
	ctx context.Context,
	req *timeZoneFetcherGRPCInterface.RequestForTimeZone,
) (*timeZoneFetcherGRPCInterface.TimeZoneResponse, error) {
	tz, err := fetcher.userService.GetUserTimeZone(req.GetDate())
	if err != nil {
		return nil, err
	}

	return &timeZoneFetcherGRPCInterface.TimeZoneResponse{Timezone: tz}, nil
}
