FROM ubuntu:16.04

MAINTAINER Ivan Nemshilov

#install front
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install curl -y
RUN curl -sL https://deb.nodesource.com/setup_6.x -o nodesource_setup.sh
RUN bash nodesource_setup.sh
RUN apt-get install nodejs -y
ADD src /IT_Berries_front
WORKDIR /IT_Berries_front
RUN npm install -y
RUN npm run build

WORKDIR /
RUN ls -al
# install nginx
RUN apt-get update
RUN apt-get install -y nginx
RUN ls -al
# Remove the default Nginx configuration file
RUN rm -v /etc/nginx/nginx.conf
# Copy a configuration file
ADD confs/nginx.conf /etc/nginx/
# add nginx conf
ADD confs/default.conf /etc/nginx/conf.d/default.conf
# Append "daemon off;" to the beginning of the configuration
RUN echo "daemon off;" >> /etc/nginx/nginx.conf
# Expose ports
EXPOSE 80
# Set the default command to execute
# when creating a new container
WORKDIR /etc/nginx
CMD ["nginx"]