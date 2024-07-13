# SPDX-FileCopyrightText: 2024 Andrew Pantuso <ajpantuso@gmail.com>
#
# SPDX-License-Identifier: MPL-2.0

FROM docker.io/golang:1.22.5 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0
RUN go build -v main.go

FROM registry.access.redhat.com/ubi8-micro@sha256:daaff371d4d735b43cc5108a1f87c68bb455573a921cb681a9a8b40ed6cc595d

WORKDIR /opt/dyn-radius

COPY --from=builder /usr/src/app/main dyn-radius

USER nobody

EXPOSE 51812
EXPOSE 8080
VOLUME [ "/opt/dyn-radius/config" ]

ENTRYPOINT [ "/opt/dyn-radius/dyn-radius" ]
