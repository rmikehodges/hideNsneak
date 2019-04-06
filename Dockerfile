FROM golang:alpine

LABEL creator="rmikehodges"
LABEL dockerfile="khast3x"
LABEL repository="https://github.com/rmikehodges/hideNsneak"

RUN apk update && apk add --no-cache bash \
                                     git \
                                     nano \
                                     python2 \
                                     py-pip \
                                     ansible \
                                     terraform
WORKDIR /opt/hidensneak
COPY . .
RUN ./setup.sh

ENTRYPOINT ["bash"]