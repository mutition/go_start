package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mutition/go_start/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.NewConsulRegistry(viper.GetString("consul.addr"))
	if err != nil {
		return nil, err
	}
	instanceID := GenerateInstanceID(serviceName)
	hostPort := viper.Sub(serviceName).Get("grpc-addr")
	if err := registry.Register(ctx, instanceID, serviceName, hostPort.(string)); err != nil {
		return nil, err
	}
	go func() {
		for {
			if err := registry.HealthCheck(ctx, instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s error: %v", instanceID, err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	logrus.WithFields(logrus.Fields{
		"instanceID":  instanceID,
		"serviceName": serviceName,
		"hostPort":    hostPort,
	}).Info("registered to consul")
	return func() error {
		return registry.DeRegister(ctx, instanceID, serviceName)
	}, nil
}

func GetServiceGRPCAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.NewConsulRegistry(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	instances, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(instances) == 0 {
		return "", fmt.Errorf("no instances found for service %s", serviceName)
	}
	i := rand.Intn(len(instances))
	logrus.Infof("Discovered %d instances for service %s, selected instance %s", len(instances), serviceName, instances[i])
	return instances[i], nil
}
