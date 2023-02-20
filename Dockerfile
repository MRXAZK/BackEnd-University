FROM golang:1.20-alpine

WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

RUN go build -o /backend-university

EXPOSE 8000

CMD [ "/backend-university" ]
