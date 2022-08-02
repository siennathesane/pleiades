FROM golang:1.18

EXPOSE 8080
RUN go install -v go.uber.org/sally@latest

ADD ref-configs/sally/sally.yaml /etc/sally.yaml
CMD ["sally", "-yml", "/etc/sally.yaml"]
