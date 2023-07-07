FROM golang:1.20 as build

WORKDIR /app

# look into copying the go.mod and sum only first
COPY . . 

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /launch cmd/launcher/main.go

FROM gcr.io/distroless/base-debian11 as release

WORKDIR /cmd/launcher

COPY --from=build /launch /launch
COPY --from=build /app/launcher/tmpl ../../launcher/tmpl

EXPOSE 80

ENTRYPOINT [ "/launch" ]