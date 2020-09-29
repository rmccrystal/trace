FROM ubuntu

# Install Golang and Node
RUN apt-get update
RUN apt-get install golang nodejs npm -y
RUN npm install -g yarn

# Install npm deps
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN yarn

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

# TODO: This could be optimized a bit by building the frontend and backend seperately
WORKDIR /app
COPY . .

WORKDIR /app/frontend
RUN yarn build

WORKDIR /app
RUN go build cmd/api

EXPOSE 8080

# node sets a different entrypoint, set it to default
ENTRYPOINT ["/bin/sh", "-c"]

CMD ["./api"]
