# nikt-link-proxy

| main | [![Build Status](https://jenkins.srv0.tokarch.uk/buildStatus/icon?job=mainnika%2Fnikt-link-proxy%2Fmain)](https://jenkins.srv0.tokarch.uk/job/mainnika/job/nikt-link-proxy/job/main) |
|------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

The proxy is basically just a link shortener. But ideally works automatically.

## the ideal pipeline

1. Do `GET nikt.tk/go?https%3A%2F%2Fexample.com`;
2. the proxy makes an short link automatically, e.g. you got `nikt.tk/uniq` and redirects to it by using `HTTP 302` redirect;
3. when proxy get a request to `nikt.tk/uniq` it makes `HTTP 301` redirect to `https://example.com` as the final destination.

## link analytics

> todo:

## referer set

> todo:

## roadmap

> todo:
