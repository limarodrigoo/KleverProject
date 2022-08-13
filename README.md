
# Crypto Vote System

A upvote system where you can create, upvote and downvote 
your cryptos



## Stacks

- Go lang
- MongoDB
- gRPC
- Gin
## Geting ready

First you need to have Go Lang installed in your system and docker

- Download Go [here](https://go.dev/dl/)

- Instructions to install Go [here](https://go.dev/doc/install)

- Guides to install docker engine [here](https://docs.docker.com/engine/install/)

Once you get all done you have to clone the repo

```bash
  git clone git@github.com:limarodrigoo/KleverProject.git
  cd KleverProject 
```

Now that you are inside the project repo directory you need to start your mongoDB with docker

```bash
docker-compose up
```

So in two terminals run the client and the server

```bash
go run server/main.go
go run client/main.go
```



## API Reference

#### Get all cryptos

```http
  GET localhost:8080/cryptos
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `-` | `-` | `-` |

List all cryptos

#### Get Crypto

```http
  GET localhost:8080/crypto/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of crypto |

List the corresponding crypto

#### Upvote Crypto

```http
  PUT localhost:8080/upvote/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of crypto |

Upvote the corresponding crypto

#### Downvote Crypto

```http
  PUT localhost:8080/downvote/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of crypto |

Downvote the corresponding crypto

#### Create Crypto

```http
  POST localhost:8080/crypto
```

| JSON Body Request | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Name`      | `string` | **Required**. name of crypto |
| `Upvote`      | `number` | **Required**. number of upvotes |
| `Downvote`      | `number` | **Required**. number of downvotes |



Create crypto
