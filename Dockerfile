FROM golang:1.26-trixie AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=1 go build -ldflags "-w -s" -o /aqi-notify

FROM gcr.io/distroless/base-debian13:nonroot AS build-release-stage

COPY --from=build-stage /aqi-notify /

ENTRYPOINT ["/aqi-notify"]
