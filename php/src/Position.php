<?php

namespace KurrentDB;

class Position
{
    public function __construct(
        public int $commit, $prepare
    )
    {

    }
}