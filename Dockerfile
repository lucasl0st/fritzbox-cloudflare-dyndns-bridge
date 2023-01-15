FROM golang:1.19-alpine

ENV GIN_MODE=release
ENV PORT=8000

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /fritzbox-cloudflare-dyndns-bridge

EXPOSE $PORT

CMD [ "/fritzbox-cloudflare-dyndns-bridge" ]