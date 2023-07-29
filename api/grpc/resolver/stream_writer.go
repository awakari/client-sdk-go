package resolver

import (
	"github.com/awakari/client-sdk-go/model"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"io"
)

type streamWriter struct {
	stream Service_SubmitMessagesClient
}

func newStreamWriter(stream Service_SubmitMessagesClient) model.Writer[*pb.CloudEvent] {
	return streamWriter{
		stream: stream,
	}
}

func (w streamWriter) Close() (err error) {
	err = w.stream.CloseSend()
	if err != nil {
		err = decodeError(err)
	}
	return
}

func (w streamWriter) WriteBatch(msgs []*pb.CloudEvent) (ackCount uint32, err error) {
	req := SubmitMessagesRequest{
		Msgs: msgs,
	}
	err = w.stream.Send(&req)
	var resp *SubmitMessagesResponse
	// when err is EOF we need to know why this happened and there's no way than receive the actual error causing this
	if err == nil || err == io.EOF {
		resp, err = w.stream.Recv()
	}
	if err != nil {
		err = decodeError(err)
	}
	if err == nil {
		ackCount = resp.AckCount
	}
	return
}
