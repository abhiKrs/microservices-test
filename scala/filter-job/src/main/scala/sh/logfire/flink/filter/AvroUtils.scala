package sh.logfire.flink.filter

import io.confluent.kafka.schemaregistry.client.CachedSchemaRegistryClient
import org.apache.avro.Schema

object AvroUtils {
  private val schemaRegistry = new CachedSchemaRegistryClient(Configs.SCHEMA_REGISTRY_ADDRESS, 1000)

  def getValueSchema(topic: String): Schema = {
    val valueSchemaString = schemaRegistry
      .getLatestSchemaMetadata(topic + "-value")
      .getSchema
    new Schema.Parser().parse(valueSchemaString)
  }
}
