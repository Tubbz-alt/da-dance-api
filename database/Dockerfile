FROM postgres:11
ENV POSTGRES_DB "dda"
ENV POSTGRES_USER "secret_user"
ENV POSTGRES_PASSWORD "secret_password"

COPY sql/ /docker-entrypoint-initdb.d/