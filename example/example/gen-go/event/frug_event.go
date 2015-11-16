// Autogenerated by Frugal Compiler (0.0.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package event

import (
	"fmt"
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal-go"
)

const delimiter = "."

type EventsPublisher struct {
	Transport frugal.Transport
	Protocol  thrift.TProtocol
	SeqId     int32
}

func NewEventsPublisher(provider *frugal.Provider) *EventsPublisher {
	transport, protocol := provider.New()
	return &EventsPublisher{
		Transport: transport,
		Protocol:  protocol,
		SeqId:     0,
	}
}

func (l *EventsPublisher) PublishEventCreated(user string, req *Event) error {
	op := "EventCreated"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	l.Transport.PreparePublish(topic)
	oprot := l.Protocol
	l.SeqId++
	if err := oprot.WriteMessageBegin(op, thrift.CALL, l.SeqId); err != nil {
		return err
	}
	if err := req.Write(oprot); err != nil {
		return err
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	return oprot.Flush()
}

type EventsSubscriber struct {
	Provider *frugal.Provider
}

func NewEventsSubscriber(provider *frugal.Provider) *EventsSubscriber {
	return &EventsSubscriber{Provider: provider}
}

func (l *EventsSubscriber) SubscribeEventCreated(user string, handler func(*Event)) (*frugal.Subscription, error) {
	op := "EventCreated"
	prefix := fmt.Sprintf("foo.%s.", user)
	topic := fmt.Sprintf("%sEvents%s%s", prefix, delimiter, op)
	transport, protocol := l.Provider.New()
	if err := transport.Subscribe(topic); err != nil {
		return nil, err
	}

	sub := frugal.NewSubscription(topic, transport)
	go func() {
		for {
			received, err := l.recvEventCreated(op, protocol)
			if err != nil {
				if e, ok := err.(thrift.TTransportException); ok && e.TypeId() == thrift.END_OF_FILE {
					return
				}
				log.Println("frugal: error receiving:", err)
				sub.Signal(err)
				sub.Unsubscribe()
				return
			}
			handler(received)
		}
	}()

	return sub, nil
}

func (l *EventsSubscriber) recvEventCreated(op string, iprot thrift.TProtocol) (*Event, error) {
	name, _, _, err := iprot.ReadMessageBegin()
	if err != nil {
		return nil, err
	}
	if name != op {
		iprot.Skip(thrift.STRUCT)
		iprot.ReadMessageEnd()
		x9 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
		return nil, x9
	}
	req := &Event{}
	if err := req.Read(iprot); err != nil {
		return nil, err
	}

	iprot.ReadMessageEnd()
	return req, nil
}

type FrugalFoo interface {
	Blah(frugal.Context, int32) (int64, error)
}

type FrugalFooClient struct {
	protocol frugal.FProtocol
	wrapped  Foo
}

func NewFrugalFooClientFactory(t thrift.TTransport, f frugal.FProtocolFactory) *FrugalFooClient {
	return &FrugalFooClient{
		protocol: f.GetProtocol(t),
		wrapped:  NewFooClientProtocol(t, f.GetProtocol(t), f.GetProtocol(t)),
	}
}

func NewFrugalFooClient(t thrift.TTransport, proto frugal.FProtocol) *FrugalFooClient {
	return &FrugalFooClient{
		protocol: proto,
		wrapped:  NewFooClientProtocol(t, proto, proto),
	}
}

func (f *FrugalFooClient) Blah(ctx frugal.Context, num int32) (r int64, err error) {
	if err = f.protocol.WriteRequestHeader(ctx); err != nil {
		return
	}
	r, err = f.wrapped.Blah(num)
	if e := f.protocol.ReadResponseHeader(ctx); e != nil {
		err = e
	}
	return r, err
}

type FrugalFooProcessor struct {
	processorMap map[string]frugal.FProcessorFunction
	handler      FrugalFoo
}

func (p *FrugalFooProcessor) GetProcessorFunction(key string) (processor frugal.FProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return
}

func NewFrugalFooProcessor(handler FrugalFoo) *FrugalFooProcessor {
	p := &FrugalFooProcessor{
		handler:      handler,
		processorMap: make(map[string]frugal.FProcessorFunction),
	}
	p.processorMap["blah"] = &fooFrugalProcessorBlah{handler: handler}
	return p
}

func (p *FrugalFooProcessor) Process(iprot, oprot frugal.FProtocol) (success bool, err thrift.TException) {
	ctx, err := iprot.ReadRequestHeader()
	if err != nil {
		return false, err
	}
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x3.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x3
}

type fooFrugalProcessorBlah struct {
	handler FrugalFoo
}

func (p *fooFrugalProcessorBlah) Process(ctx frugal.Context, seqId int32, iprot, oprot frugal.FProtocol) (success bool, err thrift.TException) {
	args := FooBlahArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("blah", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := FooBlahResult{}
	var retval int64
	var err2 error
	if retval, err2 = p.handler.Blah(ctx, args.Num); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing blah: "+err2.Error())
		oprot.WriteMessageBegin("blah", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = &retval
	}
	if err2 = oprot.WriteMessageBegin("blah", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

