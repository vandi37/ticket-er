FROM gomicro/goose AS goose

FROM scratch
COPY --from=goose /bin/goose /bin/goose

WORKDIR /migrations
COPY . /migrations