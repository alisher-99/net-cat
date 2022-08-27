FROM golang:1.17-alpine AS build 

WORKDIR /app 
COPY . .
RUN go build -o server .

FROM alpine 
WORKDIR /app
COPY --from=build /app .
CMD ./server