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

# migrate is a base container with migrate and the Cloud SQL proxy installed.

FROM golang:1.14-alpine AS builder

# Install deps
RUN apk add --no-cache git

# Install migrate
RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

COPY . .

FROM alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /go/migrations /migrations


ENTRYPOINT ["/usr/local/bin/migrate"]
