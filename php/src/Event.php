<?php

namespace KurrentDB;

use Attribute;

#[Attribute(Attribute::TARGET_CLASS)]
final class Event
{
    public function __construct(
        public string $type,
        public int    $version = 1
    )
    {
    }
}