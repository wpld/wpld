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
	    try_files $uri =404;
	    fastcgi_split_path_info ^(.+\.php)(/.+)$;

	    include /etc/nginx/fastcgi_params;
	    # Long timeout to allow more time with Xdebug
	    fastcgi_read_timeout 3600s;
	    fastcgi_buffer_size 128k;
	    fastcgi_buffers 4 128k;
	    fastcgi_pass ${PHPFPM_HOST}:9000;
	    fastcgi_index index.php;
	    fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
	}
}
`
