FROM golang:latest as build
WORKDIR /app
COPY . .
RUN go build -o out cmd/reindexerapp/main.go

FROM gcr.io/distroless/base-debian11
COPY --from=build /app/out .
EXPOSE 8000
CMD [ "out" ]
