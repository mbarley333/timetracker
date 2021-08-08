FROM golang:1.16-alpine AS build

ADD . /go/src/timetracker
WORKDIR /go/src/timetracker/cmd
RUN go test
RUN CGO_ENABLED=0 go build -o /bin/timetracker

FROM scratch
COPY --from=build /bin/timetracker /bin/timetracker
ENTRYPOINT ["/bin/timetracker"]