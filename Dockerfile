FROM golang:1.13
RUN go get -u github.com/tcnksm/ghr
ENTRYPOINT ghr -u spaceapegames $VERSION pkg/
