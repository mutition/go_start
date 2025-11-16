package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type ConsulRegistry struct {
	client *api.Client
}

var (
	consulClient *ConsulRegistry
	once         sync.Once
	initErr      error
)

func NewConsulRegistry(consulAddr string) (*ConsulRegistry, error) {
	once.Do(func() {
		config := api.DefaultConfig()
		config.Address = consulAddr
		client, err := api.NewClient(config)
		if err != nil {
			initErr = err
			return
		}
		consulClient = &ConsulRegistry{client: client}
	})
	if initErr != nil {
		return nil, initErr
	}
	return consulClient, nil
}

// implement Registry interface
func (r *ConsulRegistry) Register(_ context.Context, instanceID, serviceName, hostPort string) error {
	parseHostPort := strings.Split(hostPort, ":")
	if len(parseHostPort) != 2 {
		return errors.New("invalid host port")
	}
	host := parseHostPort[0]
	port, err := strconv.Atoi(parseHostPort[1])
	if err != nil {
		return errors.New("invalid port")
	}
	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Check: &api.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  false,
			TTL:                            "5s",
			DeregisterCriticalServiceAfter: "10s",
			Timeout:                        "5s",
		},
	})
}

func (r *ConsulRegistry) DeRegister(ctx context.Context, instanceID, serviceName string) error {
	logrus.WithFields(logrus.Fields{
		"instanceID":  instanceID,
		"serviceName": serviceName,
	}).Info("deregistering service")

	// ✅ 先注销 Service（会自动移除关联的 Check）
	if err := r.client.Agent().ServiceDeregister(instanceID); err != nil {
		logrus.Errorf("failed to deregister service: %v", err)
		return err
	}


	return nil
}

func (r *ConsulRegistry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, entry := range entries {
		res = append(res, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}
	return res, nil
}

func (r *ConsulRegistry) HealthCheck(ctx context.Context, instanceID, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
