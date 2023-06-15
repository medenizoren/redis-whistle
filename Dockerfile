FROM scratch

COPY redis-whistle /app

ENTRYPOINT [ "/app" ]