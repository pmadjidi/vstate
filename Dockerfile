FROM iron/go:dev
WORKDIR /app
ENV SRC_DIR=/go/src/vehicles
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR;go get . ; go build -o vehicle; cp vehicle /app/
ENTRYPOINT ["./vehicle"]
