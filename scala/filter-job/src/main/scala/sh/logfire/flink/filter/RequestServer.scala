package sh.logfire.flink.filter

import io.grpc.protobuf.services.ProtoReflectionService
import io.grpc.stub.StreamObserver
import io.grpc.{Server, ServerBuilder}
import org.slf4j.{Logger, LoggerFactory}
import sh.logfire.request.{FilterRequest, FilterResponse, FlinkServiceGrpc}

import scala.concurrent.{ExecutionContext, Future}

object RequestServer {
  private val logger: Logger = LoggerFactory.getLogger(getClass.getSimpleName)


  def main(args: Array[String]): Unit = {

    val server = new RequestServer(ExecutionContext.global)
    server.start()
    server.blockUntilShutdown()
  }

  private val port = 50051

}

class RequestServer(executionContext: ExecutionContext) { self =>
  private[this] var server: Server = null

  private def start(): Unit = {
    server = ServerBuilder
      .forPort(RequestServer.port)
      .addService(FlinkServiceGrpc.bindService(new FlinkServiceImpl, executionContext))
      .addService(ProtoReflectionService.newInstance())
      .build
      .start
    RequestServer.logger.info("Server started, listening on " + RequestServer.port)
    sys.addShutdownHook {
      System.err.println("*** shutting down gRPC server since JVM is shutting down")
      self.stop()
      System.err.println("*** server shut down")
    }
  }

  private def stop(): Unit =
    if (server != null) {
      server.shutdown()
    }

  private def blockUntilShutdown(): Unit =
    if (server != null) {
      server.awaitTermination()
    }

  private class FlinkServiceImpl extends FlinkServiceGrpc.FlinkService {
    override def submitFilterRequest(request: FilterRequest, responseObserver: StreamObserver[FilterResponse]): Unit = ???
  }
}
