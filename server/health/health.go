// health package implements health of the service.
//
// TODO: inject outer dependencies and return health based on ping checks. Example: if redis is
// down, the svc is not fully healthy.

package health

import (
	"context"
	"time"

	pb "github.com/toncek345/userservice/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type HealthServer struct {
	pb.UnimplementedHealthServer
}

func (h *HealthServer) Check(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (h *HealthServer) Watch(_ *emptypb.Empty, srv pb.Health_WatchServer) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		if err := srv.Send(&emptypb.Empty{}); err != nil {
			return nil
		}
		<-ticker.C
	}
}
