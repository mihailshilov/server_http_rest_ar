FROM golang:latest
RUN apt-get -y update && apt-get -y upgrade
RUN mkdir /server
ADD . /server/ 
WORKDIR /server
RUN make
CMD ["./apiserver"]
RUN find . ! -name 'apiserver' -type f -exec rm -f {} +
RUN rm -R -- */




