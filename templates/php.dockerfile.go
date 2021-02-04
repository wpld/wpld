package templates

import "github.com/MakeNowJust/heredoc"

var PHP_DOCKERFILE = heredoc.Doc(`
	ARG PHP_IMAGE=php:7.3-fpm-alpine
	FROM ${PHP_IMAGE}

	ARG CALLING_USER=www-data
	ARG CALLING_UID=82
	
	USER root
	
	RUN adduser -u ${CALLING_UID} -D -S -G www-data ${CALLING_USER}
	#RUN mkdir -p /run/php-fpm
	#RUN chown ${CALLING_USER} /run/php-fpm
	#RUN chown ${CALLING_USER} /var/log/php-fpm
	#RUN touch /usr/local/etc/msmtprc && chown ${CALLING_USER} $_
	
	USER ${CALLING_USER}
`)
