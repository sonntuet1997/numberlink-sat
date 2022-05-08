# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.17.2-alpine as base
RUN apk add --no-cache ca-certificates git curl build-base

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.5 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe 
WORKDIR /app

FROM golang:1.17.2-alpine as cadical
RUN apk add --no-cache ca-certificates git curl build-base g++ clang
WORKDIR /app
RUN git clone https://github.com/arminbiere/cadical
RUN cd cadical && ./configure && make
WORKDIR /app


FROM base as dev
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && mv ./bin/air /bin/air
RUN go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest


FROM base as builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build-server


FROM alpine as release
WORKDIR /app
RUN apk add g++ clang
COPY --from=cadical /app/cadical/build/cadical /bin/cadical
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=builder /app/bin/server /usr/bin/server
ENTRYPOINT ["/usr/bin/server"]