<?php

namespace KurrentDB;

use JsonSerializable;

enum ContentType: int
{
    case Binary = 0;
    case Json = 1;
}