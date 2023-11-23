FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /main


FROM scratch
WORKDIR /
COPY --from=build /main /main
EXPOSE 3001
ENTRYPOINT ["/main"]