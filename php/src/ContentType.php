<?php

namespace KurrentDB;

use JsonSerializable;

enum ContentType: int implements JsonSerializable
{
    case Binary = 0;
    case Json = 1;

    public function jsonSerialize(): int
    {
        return $this->value;
    }
}