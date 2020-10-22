######## Builder ########
FROM golang:1.14 as builder
WORKDIR /app
COPY / .
ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN go build -ldflags="-X config.versionTag=$(git describe --tags)" -mod=vendor -v -o app main.go


######### Container #########
FROM alpine:3
RUN apk update && apk upgrade && \
    apk add --no-cache bash curl mysql-client redis
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /app/app /app/app

# Create appuser.
ENV USER=appuser
ENV UID=54878
RUN adduser \
--disabled-password \
--gecos "application user" \
--no-create-home \
--uid "${UID}" \
"${USER}"

RUN chown -R appuser:appuser /app
USER appuser:appuser

CMD ["sleep", "100000s"]