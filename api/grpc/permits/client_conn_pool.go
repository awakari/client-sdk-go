package permits

import (
	"context"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type clientConnPool struct {
	connPool *grpcpool.Pool
}

func NewClientConnPool(connPool *grpcpool.Pool) ServiceClient {
	return clientConnPool{
		connPool: connPool,
	}
}

func (c clientConnPool) GetUsage(ctx context.Context, req *GetUsageRequest, opts ...grpc.CallOption) (resp *GetUsageResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.GetUsage(ctx, req, opts...)
	}
	return
}
