<?php

/**
 * This is needed for cookie based authentication to encrypt password in
 * cookie. Needs to be 32 chars long.
 */
$cfg['blowfish_secret'] = 'l3+wF5o$MUK@hj;[HLkQ4#V9-m?b4JmgXa]H_{uH#H]x|oQI%c1s|wFOGTc[<{3M';

/**
 * Servers configuration
 */
$i                    = 0;
$cfg['ServerDefault'] = 1;

{{ range $key, $value := . }}
$i++;
$cfg['Servers'][$i]['host']      = '{{ $key }}';
$cfg['Servers'][$i]['auth_type'] = 'config';
$cfg['Servers'][$i]['user']      = 'root';
$cfg['Servers'][$i]['password']  = 'password';
$cfg['Servers'][$i]['verbose']   = '{{ $value }}';
{{ end }}
