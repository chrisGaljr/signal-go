FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY template/ /app/template/
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    chromium \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ttf-freefont && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/template ./template

RUN chown appuser:appgroup main
RUN chown -R appuser:appgroup /root/

USER appuser

EXPOSE 7777

CMD ["./main"]