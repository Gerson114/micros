# 1️⃣ Imagem base com Go 1.25+
FROM golang:1.25

# 2️⃣ Diretório de trabalho dentro do container
WORKDIR /app

# 3️⃣ Copia go.mod e go.sum primeiro (mais eficiente)
COPY go.mod go.sum ./

# 4️⃣ Baixa dependências do Go
RUN go mod download

# 5️⃣ Copia todo o restante do código da aplicação
COPY . .

# 6️⃣ Builda o executável
RUN go build -o app

# 7️⃣ Expõe a porta que a API escuta
EXPOSE 8080

# 8️⃣ Comando padrão para iniciar a API
CMD ["./app"]