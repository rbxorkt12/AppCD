FROM golang:latest

MAINTAINER Gyutae Park (gtp7473@sk.com)

RUN mkdir /app
WORKDIR /app

COPY cmd /app/appcd

EXPOSE 8080

