FROM golang:1.23rc2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /manga-bookmarker-backend

EXPOSE 8080

CMD [ "/manga-bookmarker-backend" ]