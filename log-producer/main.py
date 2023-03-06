import argparse
import json
import logging
import time
import sys
import requests


logging.basicConfig(
    format="%(asctime)s %(name)-12s %(levelname)-8s %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
    level=logging.INFO,
    handlers=[
        #   logging.FileHandler("logs_producer.log"),
        logging.StreamHandler(sys.stdout)
    ],
)

logger = logging.getLogger()


class ProducerCallback:
    def __init__(self, record, log_success=False):
        self.record = record
        self.log_success = log_success

    def __call__(self, err, msg):
        if err:
            logger.error("Error producing record {}".format(self.record))
        elif self.log_success:
            logger.info(
                "Produced {} to topic {} partition {} offset {}".format(
                    self.record, msg.topic(), msg.partition(), msg.offset()
                )
            )


def main(args):
    logger.info("Starting logs producer")

    while True:
        with open("logs.json", "r") as f:
            data = json.load(f)

        for logs in data["data"]:
            headers = {"authorization": "bearer_1796cb38-5ad2-4f16-8804-440b7c4c12d5"}
            # headers = {"authorization": "bearer_6aba34fd-b92e-446b-8ede-38f420a17e15"}
            x = requests.post(
                # "http://localhost/api/logfire.sh",
                "https://apibeta.logfire.sh/api/logfire.sh", 
                json=logs, 
                headers=headers
            )

            print(x.status_code)
            print(x.content)
            time.sleep(0.05)

        # time.sleep(5)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--bootstrap-server", default="localhost:9092")
    parser.add_argument("--topic", default="logs-topic")
    args = parser.parse_args()
    main(args)
