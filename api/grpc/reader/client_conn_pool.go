package reader

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

func (cp clientConnPool) Read(ctx context.Context, opts ...grpc.CallOption) (stream Service_ReadClient, err error) {
	var conn *grpcpool.ClientConn
	conn, err = cp.connPool.Get(ctx)
	var c *grpc.ClientConn
	if err == nil {
		c = conn.ClientConn
		conn.Close() // return back to the conn pool immediately
	}
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(c)
		stream, err = client.Read(ctx, opts...)
	}
	return
}
