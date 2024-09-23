package limits

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

func (c clientConnPool) Get(ctx context.Context, req *GetRequest, opts ...grpc.CallOption) (resp *GetResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.Get(ctx, req, opts...)
	}
	return
}
