FROM golang:1.16-alpine AS build
WORKDIR /build
COPY . .
RUN go mod vendor && CGO_ENABLED=no go build -tags netgo -ldflags '-w' -o /vault-aws-credential-helper .

FROM scratch
LABEL org.opencontainers.image.source https://github.com/the-maldridge/vault-aws-credential-helper
COPY --from=build /vault-aws-credential-helper /vault-aws-credential-helper
ENTRYPOINT ["/vault-aws-credential-helper"]
