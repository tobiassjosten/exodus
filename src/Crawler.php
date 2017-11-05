<?php

namespace Tobiassjosten\Exodus;

use GuzzleHttp\Psr7\Response;
use Symfony\Component\DomCrawler\Crawler as DomCrawler;

class Crawler
{
    private $client;
    private $resources = [];

    public function __construct($old, $new, Client $client = null)
    {
        $this->client = $client ?? new Client();
        $this->oldOrigin = $this->origin($old);
        $this->newOrigin = $this->origin($new);
        $this->resources[] = new Resource('/');
    }

    private function origin($url)
    {
        if (!($parsed = parse_url($url)) || !isset($parsed['host'])) {
            throw new \InvalidArgumentException("Invalid website URL '$url'");
        }

        // @todo Make initial request to check for www/HTTPS redirections.

        return sprintf('%s://%s', $parsed['scheme'] ?? 'http', $parsed['host']);
    }

    public function crawl(Resource $resource = null)
    {
        if (!$resource && !($resource = $this->resource())) {
            return;
        }

        $this->extract($resource);

        $resource->response($this->client->get($resource->url($this->newOrigin)));

        return $resource;
    }

    private function resource($path = null)
    {
        foreach ($this->resources as $resource) {
            if ($path === $resource->url() || (!$path && !$resource->response())) {
                return $resource;
            }
        }
    }

    private function extract($resource)
    {
        $response = $this->client->get($resource->url($this->oldOrigin));
        if ('text/html' === trim(explode(';', $response->getHeaderLine('Content-Type'))[0])) {
            foreach ($this->paths($response->getBody()->getContents(), $this->oldOrigin) as $path) {
                if ($newResorce = $this->resource($path)) {
                    $newResorce->referrals($resource);
                } else {
                    $this->resources[] = new Resource($path, $resource);
                }
            }
        }
    }

    public function paths($contents, $origin)
    {
        $paths = [];

        $crawler = new DomCrawler($contents);

        $crawler->filterXPath('//a')->each(function (DomCrawler $node) use (&$paths, $origin) {
            if (!($parsed = parse_url($node->attr('href')))) {
                return;
            }

            if (isset($parsed['host']) && $origin !== ($parsed['scheme'] ?? 'http').'://'.$parsed['host']) {
                return;
            }

            if ('http' !== substr($parsed['scheme'] ?? 'http', 0, 4)) {
                return;
            }

            if (!empty($parsed['path'])) {
                // Skip Cloudflare protected mailto URLs.
                if ('/cdn-cgi/l/email-protection' === substr($parsed['path'], 0, 27)) {
                    return;
                }

                $paths[] = $parsed['path'];
            }
        });

        $crawler->filterXPath('//img')->each(function (DomCrawler $node) use (&$paths, $origin) {
            if (!($parsed = parse_url($node->attr('src')))) {
                return;
            }

            if (isset($parsed['host']) && $origin !== ($parsed['scheme'] ?? 'http').'://'.$parsed['host']) {
                return;
            }

            if (!empty($parsed['path'])) {
                $paths[] = $parsed['path'];
            }
        });

        return $paths;
    }
}
