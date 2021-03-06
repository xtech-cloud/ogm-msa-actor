package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"ogm-actor/config"
	"ogm-actor/handler"
	"ogm-actor/model"
	"ogm-actor/cache"
	"os"
	"path/filepath"
	"time"

	"github.com/asim/go-micro/plugins/server/grpc/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"

	proto "github.com/xtech-cloud/ogm-msp-actor/proto/actor"
)

func main() {
	config.Setup()

    model.Setup()
    defer model.Cancel()
    model.AutoMigrateDatabase()

    cache.Setup()
    defer cache.Cancel()

	// New Service
	service := micro.NewService(
        micro.Server(grpc.NewServer()),
		micro.Name(config.Schema.Service.Name),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)

	// Initialise service
	service.Init()

	// Register Handler
	proto.RegisterHealthyHandler(service.Server(), new(handler.Healthy))
	proto.RegisterDomainHandler(service.Server(), new(handler.Domain))
	proto.RegisterDeviceHandler(service.Server(), new(handler.Device))
	proto.RegisterGuardHandler(service.Server(), new(handler.Guard))
	proto.RegisterApplicationHandler(service.Server(), new(handler.Application))
	proto.RegisterSyncHandler(service.Server(), new(handler.Sync))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
