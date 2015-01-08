package main

import (
	"fmt"
	"log"
	"os"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/bohlander/gothriftexample/multiply"
)

type MyMultiplyService struct {
}

func (s *MyMultiplyService) Multiply(n1 multiply.Int, n2 multiply.Int) (r multiply.Int, err error) {

	return n1 * n2, nil
}

func RunThrift(port int, protocol string, framed, buffered bool) {
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		//Usage()
		os.Exit(1)
	}

	var transportFactory thrift.TTransportFactory
	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	var err error
	var transport thrift.TServerTransport
	transport, err = thrift.NewTServerSocket(fmt.Sprintf("%s:%d", "localhost", port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	processor := multiply.NewMultiplicationServiceProcessor(&MyMultiplyService{})

	if err == nil {
		server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
		server.Serve()
	} else {
		log.Fatal(err)
	}
}

func main() {
	RunThrift(9000, "compact", true, false)
}
