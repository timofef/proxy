FROM golang:1.15 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./main.go

FROM ubuntu:20.04

MAINTAINER Timofei Makarov

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER proxy WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -O proxy security &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/main .

EXPOSE 8081
EXPOSE 8082

ENV PGPASSWORD postgres
CMD service postgresql start &&  psql -h localhost -d security -U proxy -p 5432 -a -q -f ./database.sql && ./main