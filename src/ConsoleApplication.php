<?php

namespace Tobiassjosten\Exodus;

use Symfony\Component\Console\Application;

class ConsoleApplication extends Application
{
    public function __construct()
    {
        parent::__construct('Exodus', '1.0.0');

        $this->add($command = new CrawlCommand());
        $this->setDefaultCommand($command->getName(), true);
    }

    public function getLongVersion()
    {
        return parent::getLongVersion().' by <comment>Tobias Sjösten</comment>';
    }
}
