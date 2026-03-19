FROM nginx:stable-alpine

# Copia o arquivo de configuração do Nginx
COPY nginx.conf /etc/nginx/nginx.conf

# Porta que o Render vai expor
EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]