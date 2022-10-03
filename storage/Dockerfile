FROM alpine:latest

#RUN apt-get update && \
#     apt-get install -y ca-certificates && \
#     rm -rf /var/lib/apt/lists/*

EXPOSE 8080

RUN mkdir /app 
ADD ./storage-server /app/
WORKDIR /app 

RUN adduser -S -D -H -h /app appuser
USER appuser

ENTRYPOINT ./storage-server --port 8080 --host 0.0.0.0 --scheme http