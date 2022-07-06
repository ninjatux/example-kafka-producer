![CICD Workflof](https://github.com/ninjatux/example-kafka-producer/actions/workflows/cicd.yaml/badge.svg)

# Example Kafka Producer

Example web service that produces kafka messages when the `/produce` endpoint is called. This application was created for educational purposes.

### Local development

```
# starts the kafka cluster
docker-compose up
# starts the example web service
go run .
# test writing to kafka
curl --location --request POST 'http://127.0.0.1:8080/produce' --header 'Content-Type: text/plain' --data-raw '{"test":"yo"}'
```