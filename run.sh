#!/bin/bash

# Remove existing containers
docker container rm -f client_container
docker container rm -f server_container

# Remove the existing network
docker network rm game_network
# will prodice an error if the network does not exist, and continue

# Save the current working directory
BACK=$(pwd)

# Build client image
cd client
docker build -t client_image .
cd $BACK

# Build server image
cd server
docker build -t server_image .
cd $BACK

# Create a network for communication between containers
docker network create game_network

# Run client container
docker run -d -p 3000:3000 --network=game_network --name client_container client_image

# Run server container
docker run -d -p 8080:8080 --network=game_network --name server_container server_image
