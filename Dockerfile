FROM tutum/nginx
ADD public/ /app/
ADD sites-enabled/ /etc/nginx/sites-enabled/
CMD 'nginx'
