FROM orchardup/nginx
ADD public/ /var/www
CMD 'nginx'
