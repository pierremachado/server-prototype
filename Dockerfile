# Use uma imagem base do Go
FROM golang:1.23-alpine

# Defina o diretório de trabalho dentro do container
WORKDIR /app

# Copie todos os arquivos de código para o container
COPY . .

# Compile os quatro serviços (API, cliente, servidor A e servidor B)
RUN go build -o /bin/client_api ./api/client_api_A/api.go
RUN go build -o /bin/communication_api ./api/communication_api_A/api.go
RUN go build -o /bin/server_a ./server_A/server_A.go
RUN go build -o /bin/server_b ./server_B/server_B.go
RUN go build -o /bin/client ./client/client.go

# Definir a variável de ambiente que escolherá qual serviço rodar
# Isso pode ser passado no docker-compose ou na linha de comando `docker run`
ENV SERVICE_TYPE=client

# O comando ENTRYPOINT vai rodar com base no valor da variável SERVICE_TYPE
# Se SERVICE_TYPE for "client_api", rodará a API, se for "server_a", rodará o servidor A, etc.
ENTRYPOINT ["sh", "-c", "/bin/$SERVICE_TYPE"]
