package writer

import (
	"fmt"
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type writeStream struct {
	stream Service_SubmitMessagesClient
}

func newWriteStream(stream Service_SubmitMessagesClient) model.WriteStream[*pb.CloudEvent] {
	return writeStream{
		stream: stream,
	}
}

func (ws writeStream) Close() error {
	return ws.stream.CloseSend()
}

func (ws writeStream) WriteBatch(msgs []*pb.CloudEvent) (ackCount uint32, err error) {
	req := SubmitMessagesRequest{
		Msgs: msgs,
	}
	err = ws.stream.Send(&req)
	var resp *SubmitMessagesResponse
	if err == nil {
		resp, err = ws.stream.Recv()
	}
	err = decodeError(err)
	if resp != nil {
		ackCount = resp.AckCount
		if err == nil && resp.Err != "" {
			err = fmt.Errorf("%w: %s", ErrInternal, resp.Err)
		}
	}
	return
}
