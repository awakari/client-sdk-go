package reader

import (
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type streamAckReader struct {
	stream Service_ReadClient
}

func newStreamAckReader(stream Service_ReadClient) (r model.AckReader[[]*pb.CloudEvent]) {
	return streamAckReader{
		stream: stream,
	}
}

func (r streamAckReader) Close() error {
	return r.stream.CloseSend()
}

func (r streamAckReader) Read() (msgs []*pb.CloudEvent, err error) {
	var resp *ReadResponse
	resp, err = r.stream.Recv()
	if resp != nil {
		msgs = resp.Msgs
	}
	err = decodeError(err)
	return
}

func (r streamAckReader) Ack(count uint32) (err error) {
	reqAck := ReadRequest{
		Command: &ReadRequest_Ack{
			Ack: &ReadCommandAck{
				Count: count,
			},
		},
	}
	err = r.stream.Send(&reqAck)
	return
}
