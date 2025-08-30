<?php

namespace KurrentDB;

use JsonSerializable;

enum StreamState: string implements JsonSerializable
{

    case Any = "Any";
    case StreamExists = "StreamExists";
    case NoStream = "NoStream";

    public function jsonSerialize(): mixed
    {
        return $this->value;
    }
}
