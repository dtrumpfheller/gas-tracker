#########################
# Build
#########################

FROM golang:1.17 as builder

WORKDIR /go/src/github.com/dtrumpfheller/gas-tracker

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

COPY *.go .
COPY helpers/*.go ./helpers/
COPY gasbuddy/*.go ./gasbuddy/
COPY influxdb/*.go ./influxdb/

RUN CGO_ENABLED=0 go build -o /go/bin/app .


#########################
# Deploy
#########################

FROM gcr.io/distroless/static

USER nobody:nobody

COPY --from=builder --chown=nobody:nobody /go/bin/app /gas-tracker/

ENTRYPOINT ["/gas-tracker/app"]