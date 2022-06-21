#-------------#
# Build Stage #
#-------------#

# First pull Golang image
FROM golang:1.18-alpine as build-env

# Set environment variable
ENV APP_NAME bleve-server
ENV CMD_PATH main.go

# Copy application data into image
COPY . $GOPATH/src/bartmika/$APP_NAME
WORKDIR $GOPATH/src/bartmika/$APP_NAME

# Budild Linux 64bit application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /$APP_NAME $GOPATH/src/bartmika/$APP_NAME/$CMD_PATH

#-----------#
# Run Stage #
#-----------#

FROM alpine:3.14

# Set environment variable
ENV APP_NAME bleve-server

# Copy only required data into this image
COPY --from=build-env /$APP_NAME .

# Expose application port
EXPOSE 8001

# Start app
CMD ["./bleve-server", "serve"]

# SPECIAL THANKS:
# https://www.bacancytechnology.com/blog/dockerize-golang-application

#-------------------------------------------------------------------------------------------------------------
# HOWTO: BUILD AN IMAGE.
# docker build -t bmika/bleve-server:1.0 .

# HOWTO: RUN A CONTAINER.
# docker run -d -p 8001:8001 --name=bleve-server -e BLEVE_SERVER_ADDRESS="0.0.0.0:8001" bmika/bleve-server:1.0
#--------------------------------------------------------------------------------------------------------------
