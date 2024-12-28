# Goのベースイメージを指定
FROM golang:1.21.4
# ワーキングディレクトリを/appに設定
WORKDIR /app 
# go.modとgo.sumをワーキングディレクトリにコピー
COPY go.mod go.sum ./ 
#
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
