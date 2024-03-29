addSbtPlugin("org.jetbrains.scala" % "sbt-ide-settings" % "1.1.1")

addSbtPlugin("com.thesamet" % "sbt-protoc" % "1.0.6")

libraryDependencies ++= Seq(
  "com.thesamet.scalapb" %% "compilerplugin" % "0.11.12",
  "com.thesamet.scalapb" %% "scalapb-validate-codegen" % "0.3.3"
)

addSbtPlugin("com.eed3si9n" % "sbt-assembly" % "0.14.9")

addSbtPlugin("com.github.sbt" % "sbt-avro" % "3.4.0")
libraryDependencies += "org.apache.avro" % "avro-compiler" % "1.10.2"