FROM golang:1.19

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

COPY init_database.sql go.mod go.sum ./ 

RUN apt-get update && apt-get -y install sqlite3 && apt-get clean

COPY . .

CMD ["sh", "-c", "sqlite3 bank.db < init_database.sql ; go build ; ./fullCycle-challenge"]