package discovery

import (
	_ "common/config"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

func GenerateInstanceID(serviceName string) string {
	x := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return fmt.Sprintf("%s-%d", serviceName, x)
}

func RegisterToConsul(ctx context.Context, serviceName string) (func(), error) {
	consulServerHost := viper.GetString("consul.host")
	consulServerPort := viper.GetString("consul.port")
	client, err := NewConsulClient(fmt.Sprintf("%s:%s", consulServerHost, consulServerPort))
	if err != nil {
		return nil, err
	}

	instanceID := GenerateInstanceID(serviceName)
	clientHost := viper.Sub(serviceName).GetString("grpc-host")
	clientPort := viper.Sub(serviceName).GetString("grpc-port")
	clientAddress := fmt.Sprintf("%s:%s", clientHost, clientPort)
	if err = client.Register(ctx, instanceID, serviceName, clientAddress); err != nil {
		return nil, err
	}

	go func() {
		for {
			if err := client.HealthCheck(instanceID); err != nil {
				logrus.Panicf("no heartbeat from instance %s of service %s, err: %v", instanceID, serviceName, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return func() {
		_ = client.Deregister(ctx, instanceID, serviceName)
	}, nil
}

func GetServiceAddress(ctx context.Context, serviceName string) (string, error) {
	consulServerHost := viper.GetString("consul.host")
	consulServerPort := viper.GetString("consul.port")
	client, err := NewConsulClient(fmt.Sprintf("%s:%s", consulServerHost, consulServerPort))
	if err != nil {
		return "", err
	}

	addresses, err := client.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}

	if len(addresses) == 0 {
		return "", fmt.Errorf("no available instances found for service %s", serviceName)
	}
	logrus.Infof("%d instances of %s service found: %v", len(addresses), serviceName, addresses)

	return addresses[rand.Intn(len(addresses))], nil
}
