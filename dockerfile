FROM golang:1.20

COPY . /workspace

RUN go env -w GOPROXY=https://goproxy.cn,direct
ENV TZ=Asia/Shanghai
RUN cd /workspace && go mod tidy
RUN cd /workspace && go build *.go
EXPOSE 8080

WORKDIR /workspace

CMD ./logger