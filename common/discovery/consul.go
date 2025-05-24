package discovery

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
)

type ConsulClient struct {
	client *api.Client
}

var (
	once         sync.Once
	consulClient *ConsulClient = nil
	initErr      error         = nil
)

func NewConsulClient(address string) (*ConsulClient, error) {
	once.Do(func() {
		config := api.DefaultConfig()
		config.Address = address
		client, err := api.NewClient(config)
		if err != nil {
			initErr = err
			return
		}
		consulClient = &ConsulClient{client: client}
	})
	return consulClient, initErr
}

func (c *ConsulClient) Register(ctx context.Context, instanceID, serviceName, address string) error {
	defer logrus.WithFields(logrus.Fields{"instance_id": instanceID, "service_name": serviceName, "address": address}).Info("client registered")

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	return c.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Port:    portNum,
		Address: host,
		Check: &api.AgentServiceCheck{
			CheckID:                        instanceID,
			Timeout:                        "5s",
			TTL:                            "5s",
			TLSSkipVerify:                  false,
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (c *ConsulClient) Deregister(ctx context.Context, instanceID, serviceName string) error {
	defer logrus.WithFields(logrus.Fields{"instance_id": instanceID, "service_name": serviceName}).Info("client deregistered")
	return c.client.Agent().ServiceDeregister(instanceID)
}

func (c *ConsulClient) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(entries))
	for i, entry := range entries {
		result[i] = fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)
	}
	return result, nil
}

func (c *ConsulClient) HealthCheck(instanceID string) error {
	return c.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
