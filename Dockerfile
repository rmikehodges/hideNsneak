FROM golang:alpine

LABEL creator="rmikehodges"
LABEL dockerfile_author="khast3x"
LABEL repository="https://github.com/rmikehodges/hideNsneak"

RUN apk update && apk add --no-cache bash \
                                     git \
                                     nano \
                                     python2 \
                                     py-pip \
                                     ansible \
                                     terraform

WORKDIR /go/src/hideNsneak
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# cd terraform/backend
# terraform init -input=true
# terraform apply
# cd ../../

ENTRYPOINT [ "/bin/bash", "-c", "hideNsneak" ]
