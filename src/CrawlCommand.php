<?php

namespace Tobiassjosten\Exodus;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class CrawlCommand extends Command
{
    private $urls = [];

    protected function configure()
    {
        $this->setName('crawl')
            ->addArgument('old', InputArgument::REQUIRED, 'Old website')
            ->addArgument('new', InputArgument::REQUIRED, 'New website')
        ;
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $statuses = [];

        $crawler = new Crawler($input->getArgument('old'), $input->getArgument('new'));
        while ($resource = $crawler->crawl()) {
            $statuses[$resource->status()] = ($statuses[$resource->status()] ?? 0) + 1;
            $output->writeln($this->format($resource));
        }

        $statuses = array_map(function ($status) use ($statuses) {
            return sprintf('<%s>%d</>: %s', $this->formatStatus($status), $status, $statuses[$status]);
        }, array_keys($statuses));

        $output->writeln('');
        $output->writeln(implode('   <fg=black;options=bold>|</>   ', $statuses));
        $output->writeln('');

        return 0;
    }

    public function format(Resource $resource)
    {
        $status = $resource->status();

        $path = $resource->url();

        if (
            ($paths = $resource->response()->getHeader('X-Guzzle-Redirect-History'))
            && ($parsed = parse_url(array_pop($paths)))
            && !empty($parsed['path'])
        ) {
            if (rtrim($resource->url(), '/') === rtrim($parsed['path'], '/')) {
                $path = $parsed['path'];
            } else {
                $path .= " -> {$parsed['path']}";
            }
        }

        if (200 > $status || 400 <= $status) {
            $path .= sprintf(
                ' <fg=black;options=bold>(%s)</>',
                array_reduce($resource->referrals(), function ($carry, $referral) {
                    return ($carry ? ', ' : '').$referral->url();
                })
            );
        }

        return sprintf('<%s>%d</> %s', $this->formatStatus($status), $status, $path);
    }

    public function formatStatus($status)
    {
        return 400 <= $status ? 'fg=red;options=bold' : (300 <= $status ? 'fg=yellow' : 'fg=green');
    }
}
