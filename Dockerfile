FROM tutum/nginx
ADD public/ /app/
ADD sites-enabled/ /etc/nginx/sites-enabled/
ADD rewrite.conf /etc/nginx/sites-enabled/rewrite.conf
CMD 'nginx'
