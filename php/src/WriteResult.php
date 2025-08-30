<?php

namespace KurrentDB;

final readonly class WriteResult
{
    public function __construct(
        public int $commitPosition,
        public int $preparePosition,
        public int $nextExpectedVersion,
    )
    {
    }

    /** @param array<string,mixed> $payload */
    public static function fromRpc(array $payload): self
    {
        $cp = $payload['CommitPosition'];
        $pp = $payload['PreparePosition'];
        $ne = $payload['NextExpectedVersion'];
        return new self($cp, $pp, $ne);
    }
}
