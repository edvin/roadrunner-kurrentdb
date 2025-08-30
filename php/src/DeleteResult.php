<?php

namespace KurrentDB;

final readonly class DeleteResult
{
    public function __construct(
        public Position $position,
    )
    {
    }

    /** @param array<string,mixed> $payload */
    public static function fromRpc(array $payload): self
    {
        $position = $payload['Position'];
        return new self(new Position($position['Commit'], $position['Prepare']));
    }
}
