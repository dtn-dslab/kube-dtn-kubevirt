package kubedtn

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/digitalocean/go-openvswitch/ovs"
	common "github.com/dtn-dslab/kube-dtn-sidecar/common"
	pb "github.com/dtn-dslab/kube-dtn-sidecar/proto/v1"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type KubeDTN struct {
	pb.UnimplementedVMSidecarServer
	s           *grpc.Server
	lis         net.Listener
	redisClient *redis.Client
	ctx         context.Context
	ovsClient   *ovs.Client
	name        string
}

type Config struct {
	Port     int
	GRPCOpts []grpc.ServerOption
}

func New(cfg Config) (*KubeDTN, error) {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		return nil, fmt.Errorf("failed to get hostname")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	redisClient := common.GenerateRedisClient()

	ovsClient := ovs.New()
	// if err := common.CreateDefaultOVSBridge(ovsClient); err != nil {
	// 	return nil, fmt.Errorf("failed to create default OVS bridge: %s", err)
	// }

	m := &KubeDTN{
		s:           grpc.NewServer(cfg.GRPCOpts...),
		lis:         lis,
		ctx:         ctx,
		redisClient: redisClient,
		ovsClient:   ovsClient,
		name:        hostname,
	}

	pb.RegisterVMSidecarServer(m.s, m)
	reflection.Register(m.s)
	return m, nil
}

func (m *KubeDTN) InitStatus() error {

	return nil
}

func (m *KubeDTN) SetupStatus() error {
	mainInterfaceAddr, err := common.GetMainInterfaceAddress()
	if err != nil {
		return fmt.Errorf("failed to get main interface address: %s", err)
	}

	topoStatus, err := common.GetTopoStatusFromRedis(m.ctx, m.redisClient, m.name)
	if err != nil {
		return fmt.Errorf("failed to get topology status from redis: %s", err)
	}

	topoStatus.PodIP = mainInterfaceAddr

	statusJSON, err := json.Marshal(topoStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal topology status: %s", err)
	}

	if err := m.redisClient.Set(m.ctx, "cni_"+m.name+"_status", statusJSON, 0).Err(); err != nil {
		return fmt.Errorf("failed to set topology status to redis: %s", err)
	}
	return nil
}

func (m *KubeDTN) Serve() error {
	if err := m.InitStatus(); err != nil {
		return err
	}

	if err := m.SetupStatus(); err != nil {
		return err
	}

	return m.s.Serve(m.lis)
}

func (m *KubeDTN) GracefulStop() {
	m.s.GracefulStop()
}
