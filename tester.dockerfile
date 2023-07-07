FROM golang:1.20 as build

WORKDIR /app

# look into copying the go.mod and sum only first
COPY . . 

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /test cmd/tester/main.go

FROM gcr.io/distroless/base-debian11 as release

WORKDIR /cmd/tester

COPY --from=build /test /test
COPY --from=build /app/tester/tmpl ../../tester/tmpl

EXPOSE 80

ENTRYPOINT [ "/test" ]