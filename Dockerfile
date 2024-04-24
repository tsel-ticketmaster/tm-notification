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
FROM --platform=linux/amd64 chromedp/headless-shell:123.0.6312.86
# Install Google Chrome

WORKDIR /usr/src/app
COPY --from=builder /tmp/app .

RUN DEBIAN_FRONTEND=noninteractive TZ=Asia/Jakarta apt-get -y install tzdata

CMD ["./app"]