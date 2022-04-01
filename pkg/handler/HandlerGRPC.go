package handler

import (
	"net"

	"google.golang.org/grpc"

	"github.com/Borislavv/Translator-telegram-bot/pkg/api/grpc/service/timeZoneFetcherGRPCInterface"
	"github.com/Borislavv/Translator-telegram-bot/pkg/api/grpc/service/translatorGRPCInterface"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/loggerService"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
)

var (
	defaultGRPCProtocol = "tcp"
	defaultGRPCPort     = ":8017"
)

type HandlerGRPC struct {
	// deps.
	manager     *manager.Manager
	translator  *translator.TranslatorService
	logger      *loggerService.LoggerService
	userService *service.UserService
}

// NewHandlerGRPC - constructor of HandlerGRPC structure.
func NewHandlerGRPC(
	manager *manager.Manager,
	translator *translator.TranslatorService,
	loggerService *loggerService.LoggerService,
	userService *service.UserService,
) *HandlerGRPC {
	return &HandlerGRPC{
		manager:     manager,
		translator:  translator,
		logger:      loggerService,
		userService: userService,
	}
}

// ListenAndServe - serving gRPC server (!infinite loop: must be running in separate thread).
func (handler *HandlerGRPC) ListenAndServe() {
	protocol := handler.manager.Config.GRPCServer.Protocol
	port := handler.manager.Config.GRPCServer.Port

	if port == "" {
		port = defaultPort
	}
	if protocol == "" {
		protocol = defaultGRPCProtocol
	}

	gRPCServer := handler.createServer()

	handler.registerServices(gRPCServer)

	// infinite loop
	gRPCListner, err := net.Listen(protocol, ":"+port)
	if err != nil {
		handler.logger.Critical(err.Error())
		return
	}
	if err := gRPCServer.Serve(gRPCListner); err != nil {
		handler.logger.Critical(err.Error())
		return
	}
}

// createServer - creation of *grpc.Server structure.
func (handler *HandlerGRPC) createServer() *grpc.Server {
	return grpc.NewServer()
}

// registerServices - registration gRPC service interfaces into server.
func (handler *HandlerGRPC) registerServices(server *grpc.Server) {
	translatorGRPCInterface.RegisterTranslatorServiceServer(
		server,
		translator.NewTranslatorGRPC(handler.translator),
	)

	timeZoneFetcherGRPCInterface.RegisterTimeZoneFetcherServiceServer(
		server,
		service.NewTimeZoneFetcherGRPC(handler.userService),
	)
}
