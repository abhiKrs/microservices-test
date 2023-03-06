package sh.logfire.flink.filter

import io.confluent.kafka.schemaregistry.client.CachedSchemaRegistryClient
import org.apache.avro.Schema
import org.apache.avro.generic.GenericRecord
import org.apache.flink.api.common.typeinfo.TypeInformation
import org.apache.flink.api.java.typeutils.RowTypeInfo
import org.apache.flink.streaming.connectors.kafka.KafkaDeserializationSchema
import org.apache.flink.types.Row
import org.apache.kafka.clients.consumer.ConsumerRecord
import io.confluent.kafka.serializers.{AbstractKafkaSchemaSerDeConfig, KafkaAvroDeserializer, KafkaAvroDeserializerConfig}

import java.util

class KafkaGenericAvroDeserializationSchema(topic: String)
    extends KafkaDeserializationSchema[Row] {
  @transient var valueDeSerializer: KafkaAvroDeserializer = null
  @transient var keyDeSerializer: KafkaAvroDeserializer = null

  override val getProducedType: TypeInformation[Row] =
    AvroSerDe.getRowSchema(topic)
  def rowTypeInfo: RowTypeInfo = getProducedType.asInstanceOf[RowTypeInfo]



  val valueSchema: Schema = AvroUtils.getValueSchema(topic)



  private def checkInitialized(): Unit = {
    if (valueDeSerializer == null) {
      val props = new util.HashMap[String, String]
      props.put(
        AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG,
        Configs.SCHEMA_REGISTRY_ADDRESS
      )
      props.put(KafkaAvroDeserializerConfig.SPECIFIC_AVRO_READER_CONFIG, "false")
      val client = new CachedSchemaRegistryClient(
        Configs.SCHEMA_REGISTRY_ADDRESS,
        AbstractKafkaSchemaSerDeConfig.MAX_SCHEMAS_PER_SUBJECT_DEFAULT
      )
      valueDeSerializer = new KafkaAvroDeserializer(client, props)
      valueDeSerializer.configure(props, false)
    }

  }

  override def isEndOfStream(nextElement: Row): Boolean = false

  override def deserialize(record: ConsumerRecord[Array[Byte], Array[Byte]]): Row = {
    checkInitialized()

    if (record.value() == null) {
      return null
    }

    val kafkaTs = record.timestamp

    val valueRecord =
      valueDeSerializer.deserialize(record.topic, record.value).asInstanceOf[GenericRecord]


    AvroSerDe.convertAvroRecordToRow(valueSchema, rowTypeInfo, valueRecord)
  }
}
