<?php

namespace KurrentDB;

enum StreamPosition: string implements \JsonSerializable
{
    case Start = 'start';
    case End = 'end';

    public function jsonSerialize(): string
    {
        return $this->value;
    }
}