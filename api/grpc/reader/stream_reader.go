package reader

import (
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type streamReader struct {
	stream Service_ReadClient
}

func newStreamReader(stream Service_ReadClient) (r model.Reader[[]*pb.CloudEvent]) {
	return streamReader{
		stream: stream,
	}
}

func (r streamReader) Close() error {
	return r.stream.CloseSend()
}

func (r streamReader) Read() (msgs []*pb.CloudEvent, err error) {
	var resp *ReadResponse
	resp, err = r.stream.Recv()
	if err == nil {
		msgs = resp.Msgs
		reqAck := ReadRequest{
			Command: &ReadRequest_Ack{
				Ack: &ReadCommandAck{
					Count: uint32(len(msgs)),
				},
			},
		}
		err = r.stream.Send(&reqAck)
	}
	err = decodeError(err)
	return
}
