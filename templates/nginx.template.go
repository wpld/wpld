package templates

const NGINX_TEMPLATE = `
map $http_x_forwarded_proto $fe_https {
    default off;
    https on;
}

# Docker DNS resolver, so we can resolve things like "elasticsearch"
resolver 127.0.0.11 valid=10s;

server {
    listen 80 default_server;

    # Doesn't really matter because default server, but this way email doesn't throw errors
    server_name localhost;

    access_log   /var/log/nginx/access.log;
    error_log    /var/log/nginx/error.log;

    root /var/www/html;
    index index.php;

    if (!-e $request_filename) {
        rewrite /wp-admin$ $scheme://$host$uri/ permanent;
        rewrite ^(/[^/]+)?(/wp-.*) $2 last;
        rewrite ^(/[^/]+)?(/.*\.php) $2 last;
    }

    location /__elasticsearch/ {
        # Variable, so that it does not error when ES not included with the environment
        set $es_host http://elasticsearch:9200;

        rewrite ^/__elasticsearch/(.*) /$1 break;
        proxy_pass $es_host;
    }

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
        fastcgi_param HTTPS $fe_https;
    }

    location ~* ^.+\.(ogg|ogv|svg|svgz|eot|otf|woff|mp4|ttf|rss|atom|jpg|jpeg|gif|png|ico|zip|tgz|gz|rar|bz2|doc|xls|exe|ppt|tar|mid|midi|wav|bmp|rtf)$ {
        access_log off; log_not_found off; expires max;

        add_header Access-Control-Allow-Origin *;
        #{TRY_PROXY}
    }

    #{PROXY_URL}

    # This should match upload_max_filesize in php.ini
    client_max_body_size 150m;
}
`
