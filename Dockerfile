FROM ubuntu:22.04

VOLUME /olsrv_vol
EXPOSE 8080

RUN mkdir /app
WORKDIR /app

COPY ./bin/olsrv_bin /app
COPY ./content /app/content

CMD ["/app/olsrv_bin", "-dbn", "/olsrv_vol/olsrv_db.db"]