<?php

namespace KurrentDB;

use JsonSerializable;

enum Direction: int implements JsonSerializable
{
    case Forwards = 0;
    case Backwards = 1;

    public function jsonSerialize(): int
    {
        return $this->value;
    }
}