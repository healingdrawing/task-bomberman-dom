#!/bin/bash

# Remove existing containers
docker container rm -f frontend_container
docker container rm -f backend_container

# Remove the existing network
docker network rm social_network
# will prodice an error if the network does not exist, and continue

# Save the current working directory
BACK=$(pwd)

# Build frontend image
cd frontend
docker build -t frontend_image .
cd $BACK

# Build backend image
cd backend
docker build -t backend_image .
cd $BACK

# Create a network for communication between containers
docker network create social_network

# Run frontend container
docker run -d -p 3000:3000 --network=social_network --name frontend_container frontend_image

# Run backend container
docker run -d -p 8080:8080 --network=social_network --name backend_container backend_image
