package sh.logfire.flink.filter

import org.apache.flink.api.common.eventtime.WatermarkStrategy
import org.apache.flink.api.common.serialization.SimpleStringSchema
import org.apache.flink.connector.kafka.source.{KafkaSource, KafkaSourceBuilder}
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer
import org.apache.flink.connector.kafka.source.reader.deserializer.KafkaRecordDeserializationSchema
import org.apache.flink.streaming.api.scala.{StreamExecutionEnvironment, createTypeInformation}
import org.apache.flink.table.api.bridge.scala.StreamTableEnvironment
import org.apache.flink.types.Row

import java.util.Properties

object FilterApplication {
  def main(args: Array[String]): Unit = {
    val env: StreamExecutionEnvironment = StreamExecutionEnvironment.getExecutionEnvironment

    val tableEnv = StreamTableEnvironment.create(env)

    val brokerAddress = Configs.KAFKA_BROKERS

    val topic = "test-go-demo-logs"

    val consumerProperties = {
      val properties = new Properties()
      properties.setProperty("isolation.level", "read_committed")
      properties.setProperty("commit.offsets.on.checkpoint", "true")
      properties
    }

    val startTimeStamp: java.lang.Long = null

    val startingOffsets =
      if (startTimeStamp != null) OffsetsInitializer.timestamp(startTimeStamp) else
      OffsetsInitializer.earliest()

    val kafkaSourceBuilder: KafkaSourceBuilder[String] = KafkaSource
      .builder()
      .setBootstrapServers(brokerAddress)
      .setTopics(topic)
      .setGroupId("filter-")
      .setStartingOffsets(startingOffsets)
      .setClientIdPrefix("filter-")
      .setProperties(consumerProperties)
      .setValueOnlyDeserializer(new SimpleStringSchema())

    val source = kafkaSourceBuilder.build()

    val dataStream = env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka Source")

    env.execute()

  }

}
