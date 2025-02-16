@echo off

:: Build the Docker images
docker-compose build

:: Start the network
docker-compose up -d

:: Wait for nodes to start
timeout /t 5

:: Show logs
docker-compose logs -f 