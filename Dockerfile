FROM golang:1.17-alpine3.14 AS builder

RUN mkdir /dist /src
WORKDIR /src

# leverage layer caching
COPY ./go.* .
RUN go mod download

COPY . .

# Build release files for every target
ENV CGO_ENABLED=0
ENV GOARCH=amd64

RUN go test ./...

ENV GOOS=linux
ENV OUTDIR=/dist/${GOARCH}-${GOOS}
RUN mkdir -p ${OUTDIR}
RUN go build -o ${OUTDIR}/event-gen main.go
RUN go build -o ${OUTDIR}/tcp-server utils/tcp-server/tcp-server.go
RUN go build -o ${OUTDIR}/http-server utils/http-server/http-server.go
RUN sha256sum ${OUTDIR}/* > ${OUTDIR}/sha256sum.txt
RUN tar -C ${OUTDIR} -cvzf /dist/${GOARCH}-${GOOS}.tar.gz .

ENV GOOS=darwin
ENV OUTDIR=/dist/${GOARCH}-${GOOS}
RUN mkdir -p ${OUTDIR}
RUN go build -o ${OUTDIR}/event-gen main.go
RUN go build -o ${OUTDIR}/tcp-server utils/tcp-server/tcp-server.go
RUN go build -o ${OUTDIR}/http-server utils/http-server/http-server.go
RUN sha256sum ${OUTDIR}/* > ${OUTDIR}/sha256sum.txt
RUN tar -C ${OUTDIR} -cvzf /dist/${GOARCH}-${GOOS}.tar.gz .

ENV GOOS=windows
ENV OUTDIR=/dist/${GOARCH}-${GOOS}
RUN mkdir -p ${OUTDIR}
RUN go build -o ${OUTDIR}/event-gen.exe main.go
RUN go build -o ${OUTDIR}/tcp-server.exe utils/tcp-server/tcp-server.go
RUN go build -o ${OUTDIR}/http-server.exe utils/http-server/http-server.go
RUN sha256sum ${OUTDIR}/* > ${OUTDIR}/sha256sum.txt
RUN tar -C ${OUTDIR} -cvzf /dist/${GOARCH}-${GOOS}.tar.gz .

FROM alpine:3.14

VOLUME /export

COPY --from=builder /src /src
COPY --from=builder /dist /dist

CMD ["ash", "-c", "cp -r /dist/* /export"]
