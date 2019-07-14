FROM golang:1.12

ENV DANCESTUDIO_MYSQLADDRESS=172.13.1.2

WORKDIR /app/studio

COPY . .

RUN go install github.com/teeaa/studio/cmd/server/.

EXPOSE 8080

CMD ["server"]
