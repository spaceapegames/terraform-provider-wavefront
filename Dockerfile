FROM golang:1.8
RUN go get -u github.com/tcnksm/ghr
RUN which ghr
ENTRYPOINT ghr -u spaceapegames $VERSION ${BINARY}_${VERSION}
