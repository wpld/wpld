package templates

const NGINX_TEMPLATE = `
server {
	listen 80 default_server;
	server_name localhost;

	access_log   /var/log/nginx/access.log;
	error_log    /var/log/nginx/error.log;

	root /var/www/html;
	index index.php;

	location / {
		try_files $uri $uri/ /index.php?$args;
	}

	location ~ \.php$ {
		proxy_pass http://${PHPFPM_HOST};
		proxy_set_header Host $host;
	}
}
`
