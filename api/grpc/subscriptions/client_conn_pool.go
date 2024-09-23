package subscriptions

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

func (c clientConnPool) Create(ctx context.Context, req *CreateRequest, opts ...grpc.CallOption) (resp *CreateResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.Create(ctx, req, opts...)
	}
	return
}

func (c clientConnPool) Read(ctx context.Context, req *ReadRequest, opts ...grpc.CallOption) (resp *ReadResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.Read(ctx, req, opts...)
	}
	return
}

func (c clientConnPool) Update(ctx context.Context, req *UpdateRequest, opts ...grpc.CallOption) (resp *UpdateResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.Update(ctx, req, opts...)
	}
	return
}

func (c clientConnPool) Delete(ctx context.Context, req *DeleteRequest, opts ...grpc.CallOption) (resp *DeleteResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.Delete(ctx, req, opts...)
	}
	return
}

func (c clientConnPool) SearchOwn(ctx context.Context, req *SearchOwnRequest, opts ...grpc.CallOption) (resp *SearchOwnResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.SearchOwn(ctx, req, opts...)
	}
	return
}

func (c clientConnPool) Search(ctx context.Context, req *SearchRequest, opts ...grpc.CallOption) (resp *SearchResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = c.connPool.Get(ctx)
	defer conn.Close()
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.Search(ctx, req, opts...)
	}
	return
}
