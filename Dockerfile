FROM golang:1.19

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

COPY init-database.sh go.mod go.sum ./ 

RUN apt-get update && \
    apt install sqlite3 && \
    chmod +x init-database.sh
COPY . .

CMD ["sh", "-c", "./init-database.sh ; go build ; ./fullCycle-challenge"]