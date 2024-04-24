FROM golang:1.21.3-bullseye as builder

LABEL maintainer="patrick_m_sangian@telkomsel.co.id"

WORKDIR /usr/src/app

COPY . ./

RUN apt-get -y update \
    && apt-get -y install netcat \
        build-essential \
        openssl \
        tzdata git

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -a -o bin/app cmd/app/main.go
RUN cp bin/app /tmp/app

# Run stage
FROM --platform=linux/amd64 debian:bullseye-slim
# Install Google Chrome
RUN apt-get update && apt-get install -y wget gnupg ca-certificates \
    && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update \
    && apt-get install -y google-chrome-stable

WORKDIR /usr/src/app
COPY --from=builder /tmp/app .

RUN apt-get install -y tzdata

CMD ["./app"]