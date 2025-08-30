<?php

namespace KurrentDB;

enum StreamPosition: string implements \JsonSerializable
{
    case Start = 'Start';
    case End = 'End';

    public function jsonSerialize(): string
    {
        return $this->value;
    }
}