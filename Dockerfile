FROM golang:1.21 as compiler
WORKDIR /src/app
COPY go.mod go.sum .env credentials.json token.json ./
RUN go mod download 
COPY . .

RUN CGO_ENABLED=0 go build -o ./a.out .

FROM alpine 
COPY --from=compiler /src/app/a.out /server
COPY --from=compiler /src/app/.env src/app/credentials.json src/app/token.json /

ENTRYPOINT [ "/server" ]