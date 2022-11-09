# Set up base image
FROM golang:1.15-alpine as build

RUN mkdir /app

ADD . /app

# Set up workdir
WORKDIR /app

# Copy and download dependencies
COPY go.mod ./
COPY go.sum ./
COPY gamefiles ./gamefiles
COPY home ./home

RUN go mod download

# Build game service
RUN go build -o ./multiplayer-snake ./main.go

# Actually running the app
FROM alpine

ENV PORT=8080

COPY --from=build /app/multiplayer-snake /multiplayer-snake
COPY --from=build /app/gamefiles /gamefiles
COPY --from=build /app/home /home

EXPOSE 8080

CMD [ "/multiplayer-snake" ]
