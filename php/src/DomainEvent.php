<?php

namespace KurrentDB;

interface DomainEvent
{
    public static function fromArray(array $data): self;
}