#!/bin/sh

echo "tearing down dev setup."

echo "stop server container"
docker container stop teacherportal

echo "stop db container"
docker container stop mysql8

echo "remove server container"
docker container rm teacherportal

echo "remove db container"
docker container rm mysql8

echo "remove docker network"
docker network rm teacherportal-network

echo "remove docker image"
docker image rm teacherportal

echo "dev tear down completed."
