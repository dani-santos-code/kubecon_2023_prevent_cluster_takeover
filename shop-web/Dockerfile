# syntax=docker/dockerfile:1

FROM golang:1.20.2
RUN apt-get update \
    && apt-get install -y wget build-essential pkg-config libjpeg-dev libpng-dev libtiff-dev libgif-dev libx11-dev \
        libmagickwand-dev imagemagick libmagickcore-dev

WORKDIR /app
COPY . .

RUN go mod download
RUN pkg-config --cflags --libs MagickWand
RUN PKG_CONFIG_PATH=/usr/local/lib/pkgconfig \
    CGO_CFLAGS_ALLOW='-Xpreprocessor' \
    CGO_CFLAGS="$(pkg-config --cflags MagickWand)" \
    CGO_LDFLAGS="$(pkg-config --libs MagickWand)" \
    && go build -o /shop-demo

RUN adduser --disabled-password --gecos "" --uid 1000 appuser
RUN chown -R appuser /app
USER 1000
CMD [ "/shop-demo" ]
