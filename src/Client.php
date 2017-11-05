<?php

namespace Tobiassjosten\Exodus;

use GuzzleHttp\Client as GuzzleClient;
use GuzzleHttp\Exception\ConnectException;

class Client extends GuzzleClient
{
    public function get($url, $options = [])
    {
        try {
            return $this->request('GET', $url, [
                'allow_redirects' => ['track_redirects' => true],
                'http_errors' => false,
            ]);
        } catch (ConnectException $e) {
            return new Response(500);
        }
    }
}
