FROM golang:1.22.4 AS build-stage
    WORKDIR /app

    COPY  go.mod go.sum ./

    COPY *.go ./

    RUN CGO_ENABLED=0 GOOS=linux go build -o /api

FROM build-stage AS run-test-stage 
    RUN go test -v ./...

FROM scratch as run-release-stage 
    WORKDIR /app

    COPY --from=build-stage /api /api

    EXPOSE 8080

    CMD [ "/api" ]
