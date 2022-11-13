#-------------#
# Build Stage #
#-------------#

# First pull Golang image
FROM golang:1.18-alpine as build-env

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
COPY . /app

# Set working directory
WORKDIR /app

# Budild Linux 64bit application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bleve-server .

#-----------#
# Run Stage #
#-----------#

FROM alpine:3.14

# Copy only binary executable into this image
COPY --from=build-env /app/bleve-server .

# Expose application port
EXPOSE 8001

# Start app
CMD ["./bleve-server", "serve"]

# SPECIAL THANKS:
# https://www.bacancytechnology.com/blog/dockerize-golang-application

#-------------------------------------------------------------------------------------------------------------
# HOWTO: BUILD AN IMAGE.
# docker build -t bartmika/bleve-server:latest .

# HOWTO: RUN A CONTAINER.
# docker run -d -p 8001:8001 --name=bleve-server -e BLEVE_SERVER_ADDRESS="0.0.0.0:8001" bartmika/bleve-server:latest
#--------------------------------------------------------------------------------------------------------------
