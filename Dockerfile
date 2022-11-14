ARG NODE_VERSION=19.0.0
FROM node:${NODE_VERSION}-slim

ARG GO_VERSION=1.19.3
ARG BUD_VERSION=main

RUN node -v

# Install basic dependencies
RUN apt-get -qq update \
  && apt-get -qq -y install curl git make gcc g++ svelte \
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
WORKDIR /app
COPY . . 
RUN bud build
EXPOSE 3000

# Run the app
ENTRYPOINT [ "./bud/app" , "--log", "debug", "--listen", "0.0.0.0:3000" ]
