<?php

namespace Tobiassjosten\Exodus;

use GuzzleHttp\Psr7\Response;

class Resource
{
    private $path, $response;
    private $referrals = [];

    public function __construct($path, Resource $referral = null)
    {
        if (empty($path[0])) {
            var_dump($path, $referral);
            die;
        }
        if ('/' !== $path[0] && $referral) {
            $path = sprintf('/%s/%s', trim($referral->url(), '/'), trim($path, '/'));
        }
        $this->path = $path;

        if ($referral) {
            $this->referrals($referral);
        }
    }

    public function url($origin = '')
    {
        return "$origin{$this->path}";
    }

    public function referrals(Resource $referral = null)
    {
        if ($referral) {
            $this->referrals[] = $referral;
        }

        return $this->referrals;
    }

    public function response(Response $response = null)
    {
        if ($response) {
            $this->response = $response;
        }

        return $this->response;
    }

    public function status()
    {
        if (!$this->response) {
            throw new \Exception('No response to extract a status code from');
        }

        $status = $this->response()->getStatusCode();

        // Error statuses are reported as is.
        if (400 <= $status) {
            return $status;
        }

        // If request has been redirected, we might want to use the initial
        // redirect and not the final request.
        if ($statuses = $this->response()->getHeader('X-Guzzle-Redirect-Status-History')) {
            $paths = $this->response()->getHeader('X-Guzzle-Redirect-History');
            $parsed = parse_url(array_pop($paths));

            // Trailing slash redirects are not considered redirects.
            if (rtrim($this->url(), '/') !== rtrim($parsed['path'], '/')) {
                return array_shift($statuses);
            }

        }

        return $status;
    }
}
