<?php

namespace KurrentDB;

use JsonSerializable;

enum StreamState: string
{

    case Any = "Any";
    case StreamExists = "StreamExists";
    case NoStream = "NoStream";
}
