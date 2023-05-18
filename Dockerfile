FROM golang:1.20

COPY . /src

RUN mkdir /photos

RUN cd /src && go build -o /onthisday-bot ./cmd/onthisday-bot/main.go

# VOLUME ["/photos"]

RUN ls /photos

CMD ["/onthisday-bot", "--path", "/photos", "--job-channel", "onthisday"]
