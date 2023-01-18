ARG NODE_VERSION=19.0.0
FROM node:${NODE_VERSION}-slim as build

ARG GO_VERSION=1.18.3
ARG BUD_VERSION=main

RUN node -v

# Install basic dependencies
RUN apt-get -qq update \
  && apt-get -qq -y install curl git make gcc g++ \
  && rm -rf /var/lib/apt/lists/*

# Install Go
RUN curl -L --output - https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -xz -C /usr/local --strip-components 1
RUN go version
ENV PATH "/root/go/bin:${PATH}"

# Install Bud
RUN git clone https://github.com/livebud/bud /bud
WORKDIR /bud
RUN git checkout $BUD_VERSION
RUN make install
RUN go install .
RUN bud version

# Build your project for production
WORKDIR /builder
COPY . . 
RUN go mod download
RUN npm install
RUN bud build
RUN ls -l
RUN ls -l bud/

FROM debian:latest
WORKDIR /app
COPY --from=build /builder/bud/app .
RUN ls -l
EXPOSE 3000

# Run the app
ENTRYPOINT [ "./app" , "--log", "debug", "--listen", "0.0.0.0:3000" ]
