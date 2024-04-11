package kubedtn

import (
	"context"

	pb "github.com/dtn-dslab/kube-dtn-sidecar/proto/v1"
)

func (m *KubeDTN) AddVMLinks(ctx context.Context, in *pb.LinksBatchQuery) (*pb.BoolResponse, error) {
	return &pb.BoolResponse{Response: true}, nil
}

func (m *KubeDTN) DeleteVMLinks(ctx context.Context, in *pb.LinksBatchQuery) (*pb.BoolResponse, error) {
	// TODO implement this
	return &pb.BoolResponse{Response: true}, nil
}

func (m *KubeDTN) UpdateVMLinks(ctx context.Context, in *pb.LinksBatchQuery) (*pb.BoolResponse, error) {
	// TODO implement this
	return &pb.BoolResponse{Response: true}, nil
}
