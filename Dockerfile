FROM golang:latest as builder
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .
RUN mkdir tech-db-homework
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tech-db-homework/main -i cmd/main.go


FROM ubuntu:19.04
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV DEBIAN_FRONTEND=noninteractive
ENV PORT 5000
ENV PGVER 11
ENV POSTGRES_HOST /var/run/postgresql/
ENV POSTGRES_PORT 5432
ENV POSTGRES_DB tech-db-forum
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres

EXPOSE $PORT

RUN apt-get update && apt-get install -y postgresql-$PGVER

USER postgres

COPY schema.sql schema.sql

RUN service postgresql start &&\
    psql -U postgres -c "ALTER USER postgres PASSWORD 'postgres';" &&\
    psql -U postgres -c 'CREATE DATABASE "tech-db-forum";' &&\
    psql -U postgres -d tech-db-forum -a -f schema.sql &&\
    service postgresql stop

#COPY config/pg_hba.conf /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "shared_buffers=256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "full_page_writes=off" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

#COPY sql/ sql/
COPY --from=builder /app/tech-db-homework/main .
CMD service postgresql start && ./main
