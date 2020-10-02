FROM ubuntu

# Install Golang and Node
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install golang nodejs npm -y

# Install npm deps
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

# TODO: This could be optimized a bit by building the frontend and backend seperately
WORKDIR /app
COPY . .

WORKDIR /app/frontend
RUN npm run build

WORKDIR /app
RUN go build ./cmd/api

EXPOSE 8080

ENV LISTEN_ADDRESS 0.0.0.0:8080
ENV MONGO_URI mongodb://localhost
ENV DATABASE_NAME dev

CMD ["./api"]
