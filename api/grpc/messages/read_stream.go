package messages

import (
	"fmt"
	"github.com/awakari/client-sdk-go/api/grpc/auth"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type readStream struct {
	stream Service_ReceiveClient
}

func newReadStream(stream Service_ReceiveClient) (rs model.ReadStream[*pb.CloudEvent]) {
	return readStream{
		stream: stream,
	}
}

func (rs readStream) Close() error {
	return rs.stream.CloseSend()
}

func (rs readStream) Read() (msg *pb.CloudEvent, err error) {
	msg, err = rs.stream.Recv()
	if err == nil {
		reqAck := ReceiveRequest{
			Command: &ReceiveRequest_Ack{
				Ack: &ReceiveCommandAck{
					Ack: true,
				},
			},
		}
		err = rs.stream.Send(&reqAck)
	}
	s, isGrpcErr := status.FromError(err)
	switch {
	case !isGrpcErr || s.Code() == codes.OK:
	case s.Code() == codes.Unauthenticated:
		err = fmt.Errorf("%w: %s", auth.ErrAuth, err)
	case s.Code() == codes.NotFound:
		err = fmt.Errorf("%w: %s", ErrNotFound, err)
	default:
		err = fmt.Errorf("%w: %s", ErrInternal, err)
	}
	return
}
