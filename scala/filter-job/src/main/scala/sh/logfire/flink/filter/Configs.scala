package sh.logfire.flink.filter

object Configs {
//  val FLINK_CLUSTER_HOST = System.getenv().getOrDefault(Constants.FLINK_CLUSTER_HOST, "logfire-flink-session-cluster-rest")
  val FLINK_CLUSTER_HOST = System.getenv().getOrDefault(Constants.FLINK_CLUSTER_HOST, "192.168.1.6")
//  val FLINK_CLUSTER_PORT = System.getenv().getOrDefault(Constants.FLINK_CLUSTER_PORT, "8081")
  val FLINK_CLUSTER_PORT = System.getenv().getOrDefault(Constants.FLINK_CLUSTER_PORT, "31822")
//  val KAFKA_BROKERS = System.getenv().getOrDefault(Constants.KAFKA_BROKERS, "")
//  val KAFKA_BROKERS = System.getenv().getOrDefault(Constants.KAFKA_BROKERS, "192.168.1.6:9092")
  val KAFKA_BROKERS = System.getenv().getOrDefault(Constants.KAFKA_BROKERS, "106.214.17.222:9092")
//  val SCHEMA_REGISTRY_ADDRESS = System.getenv().getOrDefault(Constants.SCHEMA_REGISTRY_ADDRESS, "http://192.168.1.6:31081")
//  val SCHEMA_REGISTRY_ADDRESS = System.getenv().getOrDefault(Constants.SCHEMA_REGISTRY_ADDRESS, "http://192.168.1.6:30081")
  val SCHEMA_REGISTRY_ADDRESS = System.getenv().getOrDefault(Constants.SCHEMA_REGISTRY_ADDRESS, "http://106.214.17.222:30081")
  val SEVERITY_LEVEL_FIELD = System.getenv().getOrDefault(Constants.SEVERITY_LEVEL_FIELD, "level")
  val MESSAGE_FIELD = System.getenv().getOrDefault(Constants.MESSAGE_FIELD, "message")

}
