FROM alpine as BUILD

RUN apk add go git
RUN go env -w GOPROXY=direct

ADD . .
RUN go mod download

RUN go build -o /main .

FROM alpine
COPY --from=BUILD /main /main
ENTRYPOINT ["/main"]