FROM navidonskis/nginx-php7.1:latest
RUN rm -rf /var/www
RUN git clone https://github.com/hpcslag/aws-compt-web.git /var/www
RUN cd /var/www;git pull
