FROM golang:stretch AS build-env

WORKDIR /go/src/github.com/bitdao-io/bitchain

RUN apt update
RUN apt install git -y

COPY . .

RUN make build

FROM golang:stretch

RUN apt update
RUN apt install ca-certificates jq -y

WORKDIR /root

COPY --from=build-env /go/src/github.com/bitdao-io/bitchain/build/bitchain /usr/bin/bitchain

EXPOSE 26656 26657 1317 9090 8545

CMD ["bitchain"]
