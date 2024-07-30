# Choose any golang image, just make sure it doesn't have -onbuild
FROM golang:1

RUN apt-get update && apt-get -y install libopus-dev libopusfile-dev

# Everything below is copied manually from the official -onbuild image,
# with the ONBUILD keywords removed.

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

RUN go build -o /ringbot

EXPOSE 8000

# Run
CMD ["/ringbot"]