FROM iron/go:dev
WORKDIR /app
ENV SRC_DIR=/go/src/vehicles
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN sh -c 'cd $SRC_DIR;go get . ; go test; go build -o vehicle; cp vehicle /app/'
EXPOSE 8080
ENTRYPOINT ["./vehicle"]
