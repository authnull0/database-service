#builder stage
FROM golang:1.21-alpine AS builder
#RUN apk add --no-cache postgresql16-client
ENV APPHOME=/app
WORKDIR $APPHOME
COPY . ./
RUN go mod download && go mod verify
RUN go build -o /main ./app/service.go

#final stage
FROM postgres:16.0
RUN apt-get update && apt-get install -y ca-certificates
ENV APPHOME=/app
WORKDIR $APPHOME
COPY --from=builder /main ./
COPY ./conf/config.yaml ./conf/config.yaml
RUN chmod 777 ./main
EXPOSE 9234
CMD ["./main"]

