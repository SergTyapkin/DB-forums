FROM golang:latest AS builder

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main main.go


FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

# install Postgres
ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``root`` as the password and
# then create a database `forums` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "ALTER USER postgres WITH PASSWORD 'root';" &&\
    createdb -O postgres DB-forums &&\
    /etc/init.d/postgresql stop

# Expose the PostgreSQL port
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Собранный сервер
COPY --from=builder /app /app

EXPOSE 5000
ENV PGPASSWORD root
CMD service postgresql start && psql -h localhost -d DB-forums -U postgres -p 5432 -a -q -f /app/sql/initial.sql && /app/main