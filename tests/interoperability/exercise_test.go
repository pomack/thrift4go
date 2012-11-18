package simple

import (
	"net"
	"testing"
	"thrift"
)

func TestBinaryProtocol(t *testing.T) {
	pf := thrift.NewTBinaryProtocolFactoryDefault()
	tf := thrift.NewTTransportFactory()

	addr, _ := net.ResolveTCPAddr("tcp", "localhost:8080")

	channel := thrift.NewTSocketAddr(addr)
	defer channel.Close()

	if openErr := channel.Open(); openErr != nil {
		t.Fatalf("Could not open channel due to '%q'.", openErr)
	}

	effectiveTransport := tf.GetTransport(channel)

	client := NewContainerOfEnumsTestServiceClientFactory(effectiveTransport, pf)

	var request *ContainerOfEnums = NewContainerOfEnums()

	var response *ContainerOfEnums
	var anomaly error

	if response, anomaly = client.Echo(request); anomaly != nil {
		t.Fatalf("Could not get response due to '%q'.", anomaly)
	}

	if request.First != response.First {
		t.Errorf("request.First (%q) != response.First (%q)", request.First, response.First)
	}

	if request.IsSetFirst() != response.IsSetFirst() {
		t.Errorf("request.IsSetFirst() (%q) != response.IsSetFirst() (%q).", request.IsSetFirst(), response.IsSetFirst())
	}

	if request.Second != response.Second {
		t.Errorf("request.Second (%q) != response.Second (%q).", request.Second, response.Second)
	}

	if request.IsSetSecond() != response.IsSetSecond() {
		t.Errorf("request.IsSetSecond() (%q) != response.IsSetSecond() (%q).", request.IsSetSecond(), response.IsSetSecond())
	}

	if request.Third != response.Third {
		t.Errorf("request.Third (%q) != response.Third (%q).", request.Third, response.Third)
	}

	if request.IsSetThird() != response.IsSetThird() {
		t.Errorf("request.IsSetThird() (%q) != response.IsSetThird() (%q).", request.IsSetThird(), response.IsSetThird())
	}

	if request.OptionalFourth != response.OptionalFourth {
		t.Errorf("request.OptionalFourth (%q) != response.OptionalFourth (%q).", request.OptionalFourth, response.OptionalFourth)
	}

	if request.IsSetOptionalFourth() != response.IsSetOptionalFourth() {
		t.Errorf("request.IsSetOptionalFourth() (%q) != response.IsSetOptionalFourth() (%q).", request.IsSetOptionalFourth(), response.IsSetOptionalFourth())
	}

	if request.OptionalFifth != response.OptionalFifth {
		t.Errorf("request.OptionalFifth (%q) != response.OptionalFifth (%q).", request.OptionalFifth, response.OptionalFifth)
	}

	if request.IsSetOptionalFifth() != response.IsSetOptionalFifth() {
		t.Errorf("request.IsSetOptionalFifth() (%q) != response.IsSetOptionalFifth() (%q).", request.IsSetOptionalFifth(), response.IsSetOptionalFifth())
	}

	if request.OptionalSixth != response.OptionalSixth {
		t.Errorf("request.OptionalSixth (%q) != response.OptionalSixth (%q).", request.OptionalSixth, response.OptionalSixth)
	}

	if request.IsSetOptionalSixth() != response.IsSetOptionalSixth() {
		t.Errorf("request.IsSetOptionalSixth() (%q) != response.IsSetOptionalSixth() (%q).", request.IsSetOptionalSixth(), response.IsSetOptionalSixth())
	}

	if request.DefaultSeventh != response.DefaultSeventh {
		t.Errorf("request.DefaultSeventh (%q) != response.DefaultSeventh (%q).", request.DefaultSeventh, response.DefaultSeventh)
	}

	if request.IsSetDefaultSeventh() != response.IsSetDefaultSeventh() {
		t.Errorf("request.IsSetDefaultSeventh() (%q) != response.IsSetDefaultSeventh() (%q).", request.IsSetDefaultSeventh(), response.IsSetDefaultSeventh())
	}

	if request.DefaultEighth != response.DefaultEighth {
		t.Errorf("request.DefaultEighth (%q) != response.DefaultEighth (%q).", request.DefaultEighth, response.DefaultEighth)
	}

	if request.IsSetDefaultEighth() != response.IsSetDefaultEighth() {
		t.Errorf("request.IsSetDefaultEighth() (%q) != response.IsSetDefaultEighth() (%q).", request.IsSetDefaultEighth(), response.IsSetDefaultEighth())
	}

	if request.DefaultNineth != response.DefaultNineth {
		t.Errorf("request.DefaultNineth (%q) != response.DefaultNineth (%q).", request.DefaultNineth, response.DefaultNineth)
	}

	if request.IsSetDefaultNineth() != response.IsSetDefaultNineth() {
		t.Errorf("request.IsSetDefaultNineth() (%q) != response.IsSetDefaultNineth() (%q).", request.IsSetDefaultNineth(), response.IsSetDefaultNineth())
	}

	request.First = UndefinedValues_Two
	request.Second = DefinedValues_Two
	request.Third = HeterogeneousValues_Two
	request.OptionalFourth = UndefinedValues_Three
	request.OptionalFifth = DefinedValues_Three
	request.OptionalSixth = HeterogeneousValues_Three

	if response, anomaly = client.Echo(request); anomaly != nil {
		t.Fatalf("Could not get response due to '%q'.", anomaly)
	}

	if request.First != response.First {
		t.Errorf("request.First (%q) != response.First (%q)", request.First, response.First)
	}

	if request.IsSetFirst() != response.IsSetFirst() {
		t.Errorf("request.IsSetFirst() (%q) != response.IsSetFirst() (%q).", request.IsSetFirst(), response.IsSetFirst())
	}

	if request.Second != response.Second {
		t.Errorf("request.Second (%q) != response.Second (%q).", request.Second, response.Second)
	}

	if request.IsSetSecond() != response.IsSetSecond() {
		t.Errorf("request.IsSetSecond() (%q) != response.IsSetSecond() (%q).", request.IsSetSecond(), response.IsSetSecond())
	}

	if request.Third != response.Third {
		t.Errorf("request.Third (%q) != response.Third (%q).", request.Third, response.Third)
	}

	if request.IsSetThird() != response.IsSetThird() {
		t.Errorf("request.IsSetThird() (%q) != response.IsSetThird() (%q).", request.IsSetThird(), response.IsSetThird())
	}

	if request.OptionalFourth != response.OptionalFourth {
		t.Errorf("request.OptionalFourth (%q) != response.OptionalFourth (%q).", request.OptionalFourth, response.OptionalFourth)
	}

	if request.IsSetOptionalFourth() != response.IsSetOptionalFourth() {
		t.Errorf("request.IsSetOptionalFourth() (%q) != response.IsSetOptionalFourth() (%q).", request.IsSetOptionalFourth(), response.IsSetOptionalFourth())
	}

	if request.OptionalFifth != response.OptionalFifth {
		t.Errorf("request.OptionalFifth (%q) != response.OptionalFifth (%q).", request.OptionalFifth, response.OptionalFifth)
	}

	if request.IsSetOptionalFifth() != response.IsSetOptionalFifth() {
		t.Errorf("request.IsSetOptionalFifth() (%q) != response.IsSetOptionalFifth() (%q).", request.IsSetOptionalFifth(), response.IsSetOptionalFifth())
	}

	if request.OptionalSixth != response.OptionalSixth {
		t.Errorf("request.OptionalSixth (%q) != response.OptionalSixth (%q).", request.OptionalSixth, response.OptionalSixth)
	}

	if request.IsSetOptionalSixth() != response.IsSetOptionalSixth() {
		t.Errorf("request.IsSetOptionalSixth() (%q) != response.IsSetOptionalSixth() (%q).", request.IsSetOptionalSixth(), response.IsSetOptionalSixth())
	}

	if request.DefaultSeventh != response.DefaultSeventh {
		t.Errorf("request.DefaultSeventh (%q) != response.DefaultSeventh (%q).", request.DefaultSeventh, response.DefaultSeventh)
	}

	if request.IsSetDefaultSeventh() != response.IsSetDefaultSeventh() {
		t.Errorf("request.IsSetDefaultSeventh() (%q) != response.IsSetDefaultSeventh() (%q).", request.IsSetDefaultSeventh(), response.IsSetDefaultSeventh())
	}

	if request.DefaultEighth != response.DefaultEighth {
		t.Errorf("request.DefaultEighth (%q) != response.DefaultEighth (%q).", request.DefaultEighth, response.DefaultEighth)
	}

	if request.IsSetDefaultEighth() != response.IsSetDefaultEighth() {
		t.Errorf("request.IsSetDefaultEighth() (%q) != response.IsSetDefaultEighth() (%q).", request.IsSetDefaultEighth(), response.IsSetDefaultEighth())
	}

	if request.DefaultNineth != response.DefaultNineth {
		t.Errorf("request.DefaultNineth (%q) != response.DefaultNineth (%q).", request.DefaultNineth, response.DefaultNineth)
	}

	if request.IsSetDefaultNineth() != response.IsSetDefaultNineth() {
		t.Errorf("request.IsSetDefaultNineth() (%q) != response.IsSetDefaultNineth() (%q).", request.IsSetDefaultNineth(), response.IsSetDefaultNineth())
	}
}

func TestCompactProtocol(t *testing.T) {
	pf := thrift.NewTCompactProtocolFactory()
	tf := thrift.NewTTransportFactory()

	addr, _ := net.ResolveTCPAddr("tcp", "localhost:8082")

	channel := thrift.NewTSocketAddr(addr)
	defer channel.Close()

	if openErr := channel.Open(); openErr != nil {
		t.Fatalf("Could not open channel due to '%q'.", openErr)
	}

	effectiveTransport := tf.GetTransport(channel)

	client := NewContainerOfEnumsTestServiceClientFactory(effectiveTransport, pf)

	var request *ContainerOfEnums = NewContainerOfEnums()

	var response *ContainerOfEnums
	var anomaly error

	if response, anomaly = client.Echo(request); anomaly != nil {
		t.Fatalf("Could not get response due to '%q'.", anomaly)
	}

	if request.First != response.First {
		t.Errorf("request.First (%q) != response.First (%q)", request.First, response.First)
	}

	if request.IsSetFirst() != response.IsSetFirst() {
		t.Errorf("request.IsSetFirst() (%q) != response.IsSetFirst() (%q).", request.IsSetFirst(), response.IsSetFirst())
	}

	if request.Second != response.Second {
		t.Errorf("request.Second (%q) != response.Second (%q).", request.Second, response.Second)
	}

	if request.IsSetSecond() != response.IsSetSecond() {
		t.Errorf("request.IsSetSecond() (%q) != response.IsSetSecond() (%q).", request.IsSetSecond(), response.IsSetSecond())
	}

	if request.Third != response.Third {
		t.Errorf("request.Third (%q) != response.Third (%q).", request.Third, response.Third)
	}

	if request.IsSetThird() != response.IsSetThird() {
		t.Errorf("request.IsSetThird() (%q) != response.IsSetThird() (%q).", request.IsSetThird(), response.IsSetThird())
	}

	if request.OptionalFourth != response.OptionalFourth {
		t.Errorf("request.OptionalFourth (%q) != response.OptionalFourth (%q).", request.OptionalFourth, response.OptionalFourth)
	}

	if request.IsSetOptionalFourth() != response.IsSetOptionalFourth() {
		t.Errorf("request.IsSetOptionalFourth() (%q) != response.IsSetOptionalFourth() (%q).", request.IsSetOptionalFourth(), response.IsSetOptionalFourth())
	}

	if request.OptionalFifth != response.OptionalFifth {
		t.Errorf("request.OptionalFifth (%q) != response.OptionalFifth (%q).", request.OptionalFifth, response.OptionalFifth)
	}

	if request.IsSetOptionalFifth() != response.IsSetOptionalFifth() {
		t.Errorf("request.IsSetOptionalFifth() (%q) != response.IsSetOptionalFifth() (%q).", request.IsSetOptionalFifth(), response.IsSetOptionalFifth())
	}

	if request.OptionalSixth != response.OptionalSixth {
		t.Errorf("request.OptionalSixth (%q) != response.OptionalSixth (%q).", request.OptionalSixth, response.OptionalSixth)
	}

	if request.IsSetOptionalSixth() != response.IsSetOptionalSixth() {
		t.Errorf("request.IsSetOptionalSixth() (%q) != response.IsSetOptionalSixth() (%q).", request.IsSetOptionalSixth(), response.IsSetOptionalSixth())
	}

	if request.DefaultSeventh != response.DefaultSeventh {
		t.Errorf("request.DefaultSeventh (%q) != response.DefaultSeventh (%q).", request.DefaultSeventh, response.DefaultSeventh)
	}

	if request.IsSetDefaultSeventh() != response.IsSetDefaultSeventh() {
		t.Errorf("request.IsSetDefaultSeventh() (%q) != response.IsSetDefaultSeventh() (%q).", request.IsSetDefaultSeventh(), response.IsSetDefaultSeventh())
	}

	if request.DefaultEighth != response.DefaultEighth {
		t.Errorf("request.DefaultEighth (%q) != response.DefaultEighth (%q).", request.DefaultEighth, response.DefaultEighth)
	}

	if request.IsSetDefaultEighth() != response.IsSetDefaultEighth() {
		t.Errorf("request.IsSetDefaultEighth() (%q) != response.IsSetDefaultEighth() (%q).", request.IsSetDefaultEighth(), response.IsSetDefaultEighth())
	}

	if request.DefaultNineth != response.DefaultNineth {
		t.Errorf("request.DefaultNineth (%q) != response.DefaultNineth (%q).", request.DefaultNineth, response.DefaultNineth)
	}

	if request.IsSetDefaultNineth() != response.IsSetDefaultNineth() {
		t.Errorf("request.IsSetDefaultNineth() (%q) != response.IsSetDefaultNineth() (%q).", request.IsSetDefaultNineth(), response.IsSetDefaultNineth())
	}

	request.First = UndefinedValues_Two
	request.Second = DefinedValues_Two
	request.Third = HeterogeneousValues_Two
	request.OptionalFourth = UndefinedValues_Three
	request.OptionalFifth = DefinedValues_Three
	request.OptionalSixth = HeterogeneousValues_Three

	if response, anomaly = client.Echo(request); anomaly != nil {
		t.Fatalf("Could not get response due to '%q'.", anomaly)
	}

	if request.First != response.First {
		t.Errorf("request.First (%q) != response.First (%q)", request.First, response.First)
	}

	if request.IsSetFirst() != response.IsSetFirst() {
		t.Errorf("request.IsSetFirst() (%q) != response.IsSetFirst() (%q).", request.IsSetFirst(), response.IsSetFirst())
	}

	if request.Second != response.Second {
		t.Errorf("request.Second (%q) != response.Second (%q).", request.Second, response.Second)
	}

	if request.IsSetSecond() != response.IsSetSecond() {
		t.Errorf("request.IsSetSecond() (%q) != response.IsSetSecond() (%q).", request.IsSetSecond(), response.IsSetSecond())
	}

	if request.Third != response.Third {
		t.Errorf("request.Third (%q) != response.Third (%q).", request.Third, response.Third)
	}

	if request.IsSetThird() != response.IsSetThird() {
		t.Errorf("request.IsSetThird() (%q) != response.IsSetThird() (%q).", request.IsSetThird(), response.IsSetThird())
	}

	if request.OptionalFourth != response.OptionalFourth {
		t.Errorf("request.OptionalFourth (%q) != response.OptionalFourth (%q).", request.OptionalFourth, response.OptionalFourth)
	}

	if request.IsSetOptionalFourth() != response.IsSetOptionalFourth() {
		t.Errorf("request.IsSetOptionalFourth() (%q) != response.IsSetOptionalFourth() (%q).", request.IsSetOptionalFourth(), response.IsSetOptionalFourth())
	}

	if request.OptionalFifth != response.OptionalFifth {
		t.Errorf("request.OptionalFifth (%q) != response.OptionalFifth (%q).", request.OptionalFifth, response.OptionalFifth)
	}

	if request.IsSetOptionalFifth() != response.IsSetOptionalFifth() {
		t.Errorf("request.IsSetOptionalFifth() (%q) != response.IsSetOptionalFifth() (%q).", request.IsSetOptionalFifth(), response.IsSetOptionalFifth())
	}

	if request.OptionalSixth != response.OptionalSixth {
		t.Errorf("request.OptionalSixth (%q) != response.OptionalSixth (%q).", request.OptionalSixth, response.OptionalSixth)
	}

	if request.IsSetOptionalSixth() != response.IsSetOptionalSixth() {
		t.Errorf("request.IsSetOptionalSixth() (%q) != response.IsSetOptionalSixth() (%q).", request.IsSetOptionalSixth(), response.IsSetOptionalSixth())
	}

	if request.DefaultSeventh != response.DefaultSeventh {
		t.Errorf("request.DefaultSeventh (%q) != response.DefaultSeventh (%q).", request.DefaultSeventh, response.DefaultSeventh)
	}

	if request.IsSetDefaultSeventh() != response.IsSetDefaultSeventh() {
		t.Errorf("request.IsSetDefaultSeventh() (%q) != response.IsSetDefaultSeventh() (%q).", request.IsSetDefaultSeventh(), response.IsSetDefaultSeventh())
	}

	if request.DefaultEighth != response.DefaultEighth {
		t.Errorf("request.DefaultEighth (%q) != response.DefaultEighth (%q).", request.DefaultEighth, response.DefaultEighth)
	}

	if request.IsSetDefaultEighth() != response.IsSetDefaultEighth() {
		t.Errorf("request.IsSetDefaultEighth() (%q) != response.IsSetDefaultEighth() (%q).", request.IsSetDefaultEighth(), response.IsSetDefaultEighth())
	}

	if request.DefaultNineth != response.DefaultNineth {
		t.Errorf("request.DefaultNineth (%q) != response.DefaultNineth (%q).", request.DefaultNineth, response.DefaultNineth)
	}

	if request.IsSetDefaultNineth() != response.IsSetDefaultNineth() {
		t.Errorf("request.IsSetDefaultNineth() (%q) != response.IsSetDefaultNineth() (%q).", request.IsSetDefaultNineth(), response.IsSetDefaultNineth())
	}
}

func TestJSONProtocol(t *testing.T) {
	pf := thrift.NewTJSONProtocolFactory()
	tf := thrift.NewTTransportFactory()

	addr, _ := net.ResolveTCPAddr("tcp", "localhost:8081")

	channel := thrift.NewTSocketAddr(addr)
	defer channel.Close()

	if openErr := channel.Open(); openErr != nil {
		t.Fatalf("Could not open channel due to '%q'.", openErr)
	}

	effectiveTransport := tf.GetTransport(channel)

	client := NewContainerOfEnumsTestServiceClientFactory(effectiveTransport, pf)

	var request *ContainerOfEnums = NewContainerOfEnums()

	var response *ContainerOfEnums
	var anomaly error

	if response, anomaly = client.Echo(request); anomaly != nil {
		t.Fatalf("Could not get response due to '%q'.", anomaly)
	}

	if request.First != response.First {
		t.Errorf("request.First (%q) != response.First (%q)", request.First, response.First)
	}

	if request.IsSetFirst() != response.IsSetFirst() {
		t.Errorf("request.IsSetFirst() (%q) != response.IsSetFirst() (%q).", request.IsSetFirst(), response.IsSetFirst())
	}

	if request.Second != response.Second {
		t.Errorf("request.Second (%q) != response.Second (%q).", request.Second, response.Second)
	}

	if request.IsSetSecond() != response.IsSetSecond() {
		t.Errorf("request.IsSetSecond() (%q) != response.IsSetSecond() (%q).", request.IsSetSecond(), response.IsSetSecond())
	}

	if request.Third != response.Third {
		t.Errorf("request.Third (%q) != response.Third (%q).", request.Third, response.Third)
	}

	if request.IsSetThird() != response.IsSetThird() {
		t.Errorf("request.IsSetThird() (%q) != response.IsSetThird() (%q).", request.IsSetThird(), response.IsSetThird())
	}

	if request.OptionalFourth != response.OptionalFourth {
		t.Errorf("request.OptionalFourth (%q) != response.OptionalFourth (%q).", request.OptionalFourth, response.OptionalFourth)
	}

	if request.IsSetOptionalFourth() != response.IsSetOptionalFourth() {
		t.Errorf("request.IsSetOptionalFourth() (%q) != response.IsSetOptionalFourth() (%q).", request.IsSetOptionalFourth(), response.IsSetOptionalFourth())
	}

	if request.OptionalFifth != response.OptionalFifth {
		t.Errorf("request.OptionalFifth (%q) != response.OptionalFifth (%q).", request.OptionalFifth, response.OptionalFifth)
	}

	if request.IsSetOptionalFifth() != response.IsSetOptionalFifth() {
		t.Errorf("request.IsSetOptionalFifth() (%q) != response.IsSetOptionalFifth() (%q).", request.IsSetOptionalFifth(), response.IsSetOptionalFifth())
	}

	if request.OptionalSixth != response.OptionalSixth {
		t.Errorf("request.OptionalSixth (%q) != response.OptionalSixth (%q).", request.OptionalSixth, response.OptionalSixth)
	}

	if request.IsSetOptionalSixth() != response.IsSetOptionalSixth() {
		t.Errorf("request.IsSetOptionalSixth() (%q) != response.IsSetOptionalSixth() (%q).", request.IsSetOptionalSixth(), response.IsSetOptionalSixth())
	}

	if request.DefaultSeventh != response.DefaultSeventh {
		t.Errorf("request.DefaultSeventh (%q) != response.DefaultSeventh (%q).", request.DefaultSeventh, response.DefaultSeventh)
	}

	if request.IsSetDefaultSeventh() != response.IsSetDefaultSeventh() {
		t.Errorf("request.IsSetDefaultSeventh() (%q) != response.IsSetDefaultSeventh() (%q).", request.IsSetDefaultSeventh(), response.IsSetDefaultSeventh())
	}

	if request.DefaultEighth != response.DefaultEighth {
		t.Errorf("request.DefaultEighth (%q) != response.DefaultEighth (%q).", request.DefaultEighth, response.DefaultEighth)
	}

	if request.IsSetDefaultEighth() != response.IsSetDefaultEighth() {
		t.Errorf("request.IsSetDefaultEighth() (%q) != response.IsSetDefaultEighth() (%q).", request.IsSetDefaultEighth(), response.IsSetDefaultEighth())
	}

	if request.DefaultNineth != response.DefaultNineth {
		t.Errorf("request.DefaultNineth (%q) != response.DefaultNineth (%q).", request.DefaultNineth, response.DefaultNineth)
	}

	if request.IsSetDefaultNineth() != response.IsSetDefaultNineth() {
		t.Errorf("request.IsSetDefaultNineth() (%q) != response.IsSetDefaultNineth() (%q).", request.IsSetDefaultNineth(), response.IsSetDefaultNineth())
	}

	request.First = UndefinedValues_Two
	request.Second = DefinedValues_Two
	request.Third = HeterogeneousValues_Two
	request.OptionalFourth = UndefinedValues_Three
	request.OptionalFifth = DefinedValues_Three
	request.OptionalSixth = HeterogeneousValues_Three

	if response, anomaly = client.Echo(request); anomaly != nil {
		t.Fatalf("Could not get response due to '%q'.", anomaly)
	}

	if request.First != response.First {
		t.Errorf("request.First (%q) != response.First (%q)", request.First, response.First)
	}

	if request.IsSetFirst() != response.IsSetFirst() {
		t.Errorf("request.IsSetFirst() (%q) != response.IsSetFirst() (%q).", request.IsSetFirst(), response.IsSetFirst())
	}

	if request.Second != response.Second {
		t.Errorf("request.Second (%q) != response.Second (%q).", request.Second, response.Second)
	}

	if request.IsSetSecond() != response.IsSetSecond() {
		t.Errorf("request.IsSetSecond() (%q) != response.IsSetSecond() (%q).", request.IsSetSecond(), response.IsSetSecond())
	}

	if request.Third != response.Third {
		t.Errorf("request.Third (%q) != response.Third (%q).", request.Third, response.Third)
	}

	if request.IsSetThird() != response.IsSetThird() {
		t.Errorf("request.IsSetThird() (%q) != response.IsSetThird() (%q).", request.IsSetThird(), response.IsSetThird())
	}

	if request.OptionalFourth != response.OptionalFourth {
		t.Errorf("request.OptionalFourth (%q) != response.OptionalFourth (%q).", request.OptionalFourth, response.OptionalFourth)
	}

	if request.IsSetOptionalFourth() != response.IsSetOptionalFourth() {
		t.Errorf("request.IsSetOptionalFourth() (%q) != response.IsSetOptionalFourth() (%q).", request.IsSetOptionalFourth(), response.IsSetOptionalFourth())
	}

	if request.OptionalFifth != response.OptionalFifth {
		t.Errorf("request.OptionalFifth (%q) != response.OptionalFifth (%q).", request.OptionalFifth, response.OptionalFifth)
	}

	if request.IsSetOptionalFifth() != response.IsSetOptionalFifth() {
		t.Errorf("request.IsSetOptionalFifth() (%q) != response.IsSetOptionalFifth() (%q).", request.IsSetOptionalFifth(), response.IsSetOptionalFifth())
	}

	if request.OptionalSixth != response.OptionalSixth {
		t.Errorf("request.OptionalSixth (%q) != response.OptionalSixth (%q).", request.OptionalSixth, response.OptionalSixth)
	}

	if request.IsSetOptionalSixth() != response.IsSetOptionalSixth() {
		t.Errorf("request.IsSetOptionalSixth() (%q) != response.IsSetOptionalSixth() (%q).", request.IsSetOptionalSixth(), response.IsSetOptionalSixth())
	}

	if request.DefaultSeventh != response.DefaultSeventh {
		t.Errorf("request.DefaultSeventh (%q) != response.DefaultSeventh (%q).", request.DefaultSeventh, response.DefaultSeventh)
	}

	if request.IsSetDefaultSeventh() != response.IsSetDefaultSeventh() {
		t.Errorf("request.IsSetDefaultSeventh() (%q) != response.IsSetDefaultSeventh() (%q).", request.IsSetDefaultSeventh(), response.IsSetDefaultSeventh())
	}

	if request.DefaultEighth != response.DefaultEighth {
		t.Errorf("request.DefaultEighth (%q) != response.DefaultEighth (%q).", request.DefaultEighth, response.DefaultEighth)
	}

	if request.IsSetDefaultEighth() != response.IsSetDefaultEighth() {
		t.Errorf("request.IsSetDefaultEighth() (%q) != response.IsSetDefaultEighth() (%q).", request.IsSetDefaultEighth(), response.IsSetDefaultEighth())
	}

	if request.DefaultNineth != response.DefaultNineth {
		t.Errorf("request.DefaultNineth (%q) != response.DefaultNineth (%q).", request.DefaultNineth, response.DefaultNineth)
	}

	if request.IsSetDefaultNineth() != response.IsSetDefaultNineth() {
		t.Errorf("request.IsSetDefaultNineth() (%q) != response.IsSetDefaultNineth() (%q).", request.IsSetDefaultNineth(), response.IsSetDefaultNineth())
	}
}
