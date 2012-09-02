
package thrift4go;


import org.apache.thrift.protocol.TBinaryProtocol;
import org.apache.thrift.protocol.TCompactProtocol;
import org.apache.thrift.protocol.TJSONProtocol;
import org.apache.thrift.protocol.TProtocolFactory;
import org.apache.thrift.protocol.TSimpleJSONProtocol;
import org.apache.thrift.server.TServer;
import org.apache.thrift.server.TServer.Args;
import org.apache.thrift.server.TSimpleServer;
import org.apache.thrift.transport.TServerSocket;
import org.apache.thrift.transport.TServerTransport;
import org.apache.thrift.transport.TTransportException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import thrift4go.generated.ContainerOfEnumsTestService;


public class EchoServerEntryPoint {
  private static final Logger log = LoggerFactory.getLogger(EchoServerEntryPoint.class);

  public static void main(final String[] args) throws TTransportException {
    if (args == null || args.length != 2) {
      log.warn("Expects <protocol> <port> arguments.");
      System.exit(1);
    }

    final String protocol = args[0].toUpperCase();
    final int port = Integer.parseInt(args[1]);

    log.info("Preparing to start echo service.");

    final ContainerOfEnumsTestService.Processor<EchoServiceDefinition> processor =
        new ContainerOfEnumsTestService.Processor<EchoServiceDefinition>(
            new EchoServiceDefinition(protocol));
    final TServerTransport transport = new TServerSocket(port);

    final Args serviceArguments = new Args(transport);
    serviceArguments.processor(processor);
    serviceArguments.protocolFactory(Enum.valueOf(Protocol.class, protocol).getFactory());

    final TServer server = new TSimpleServer(serviceArguments);

    log.info("Provisioned everything; now serving {} requests on {}...", protocol,port);

    try {
      server.serve();
    } finally {
      log.info("Closing down everything.");

      server.stop();
    }
  }

  private static enum Protocol {
    JSON(new TJSONProtocol.Factory()),
    SIMPLE_JSON(new TSimpleJSONProtocol.Factory()),
    BINARY(new TBinaryProtocol.Factory()),
    COMPACT(new TCompactProtocol.Factory());

    private final TProtocolFactory factory;

    Protocol(final TProtocolFactory factory) {
      this.factory = factory;
    }

    public TProtocolFactory getFactory() {
      return this.factory;
    }
  }
}
