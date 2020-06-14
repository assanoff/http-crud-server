#docker build -t crud-server .
#docker run -p 3220:8080 --name crud-server --rm -d -t crud-server

# Initial stage: download modules
FROM golang:1.13 as modules
ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.13 AS builder
COPY --from=modules /go/pkg /go/pkg
RUN mkdir /app
COPY . /app
WORKDIR /app
# We want to build our application's binary executable
RUN CGO_ENABLED=0 GOOS=linux make build

# Final stage: Run the binary
FROM scratch 
COPY --from=builder /app/bin/ .
EXPOSE 8080
CMD ["./apiserver"]