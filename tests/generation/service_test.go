package simple

import (
	"fmt"
	"net"
	"testing"
	"thrift"
	"time"
)

type ContainerOfEnumsTestService struct{}

func (c *ContainerOfEnumsTestService) Echo(message *ContainerOfEnums) (*ContainerOfEnums, error) {
	fmt.Printf("Received: %q\n", message)
	return message, nil
}

func TestGeneratedServiceInterfaceAndImplementation(t *testing.T) {
	var service IContainerOfEnumsTestService = &ContainerOfEnumsTestService{}

	input := NewContainerOfEnums()

	input.First = UndefinedTwo
	input.Second = DefinedTwo
	input.Third = HeterogeneousTwo
	input.OptionalFourth = UndefinedThree
	input.OptionalFifth = DefinedThree
	input.OptionalSixth = HeterogeneousThree
	input.DefaultSeventh = UndefinedOne
	input.DefaultEighth = DefinedOne
	input.DefaultNineth = HeterogeneousOne

	output, err := service.Echo(input)

	if err != nil {
		t.Errorf("err != nil")
	}

	if input != output {
		t.Errorf("input != output")
	}

	if *input != *output {
		t.Errorf("*input != *output")
	}
}

type ContainerOfEnumsTestServer struct {
	handler   IContainerOfEnumsTestService
	processor *ContainerOfEnumsTestServiceProcessor
}

func NewContainerOfEnumsTestServer() *ContainerOfEnumsTestServer {
	handler := &ContainerOfEnumsTestService{}
	processor := NewContainerOfEnumsTestServiceProcessor(handler)

	return &ContainerOfEnumsTestServer{
		handler:   handler,
		processor: processor,
	}
}

func TestEndToEnd(t *testing.T) {
	protocols := []thrift.TProtocolFactory{
		//		thrift.NewTCompactProtocolFactory(),
//		thrift.NewTSimpleJSONProtocolFactory(),
		thrift.NewTJSONProtocolFactory(),
		//		thrift.NewTBinaryProtocolFactoryDefault(),
	}

	for _, protocol := range protocols {
		addr, addrErr := net.ResolveTCPAddr("tcp", "localhost:9090")
		deadlineChannel := make(chan bool, 1)
		var finishedServing = false

		if addrErr != nil {
			t.Errorf("addrErr = %q", addrErr)
		}

		transport, transportErr := thrift.NewTServerSocketAddr(addr)

		if transportErr != nil {
			t.Errorf("transportErr = %q", transportErr)

			continue
		}

		processor := NewContainerOfEnumsTestServer()
		server := thrift.NewTSimpleServer4(processor.processor, transport, thrift.NewTTransportFactory(), protocol)
		completionChannel := make(chan bool, 1)

		fmt.Println("Starting server...")
		go func() {
			for !finishedServing {
				fmt.Println("Serving...")
				servingAnomaly := server.Serve()

				if servingAnomaly != nil {
					fmt.Printf("Serving Anomaly: %q", servingAnomaly)
				}
			}
		}()

		fmt.Println("Giving server some time to get ready...")
		time.Sleep(1 * time.Second)

		fmt.Println("Starting deadline timer...")
		go func() {
			time.Sleep(10 * time.Second)
			deadlineChannel <- true
		}()

		fmt.Println("Starting client...")
		go func() {
			fmt.Println("Resolving server address...")
			serverAddress, serverAddressErr := net.ResolveTCPAddr("tcp", "localhost:9090")

			if serverAddressErr != nil {
				t.Error("Could not resolve server address for client construction: %q\n", serverAddressErr)

				return
			}

			fmt.Println("Creating client transport...")
			clientTransport, clientTransportErr := thrift.NewTNonblockingSocketAddr(serverAddress)

			if clientTransportErr != nil {
				t.Errorf("Could not set up client transport: %q\n", clientTransportErr)

				return
			}

			fmt.Println("Opening transport to server...")
			if clientEstablishmentErr := clientTransport.Open(); clientEstablishmentErr != nil {
				t.Errorf("Could not establish connection with server: %q\n", clientEstablishmentErr)
			}

			fmt.Println("Preparing payload...")
			input := NewContainerOfEnums()

			input.First = UndefinedTwo
			input.Second = DefinedTwo
			input.Third = HeterogeneousTwo
			input.OptionalFourth = UndefinedThree
			input.OptionalFifth = DefinedThree
			input.OptionalSixth = HeterogeneousThree
			input.DefaultSeventh = UndefinedOne
			input.DefaultEighth = DefinedOne
			input.DefaultNineth = HeterogeneousOne

			fmt.Println("Creating client...")
			client := NewContainerOfEnumsTestServiceClientFactory(clientTransport, protocol)

			fmt.Println("Sending request...")
			emission, echoErr := client.Echo(input)

			if echoErr != nil {
				t.Errorf("Error while echoing: %q\n", echoErr)
			}

			if emission != input {
				t.Errorf("%q v. %q", emission, input)
			}

			if clientTransport != nil {
			//	clientTransport.Close()
			}

			completionChannel <- true
			finishedServing = true
		}()

		select {
		case <-completionChannel:
			fmt.Printf("Completed for %q\n", protocol)
		case <-deadlineChannel:
			t.Errorf("Deadlined for %q\n", protocol)
		}

		if transport != nil {
			// transport.Close()
		}
	}
}
