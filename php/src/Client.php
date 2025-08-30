<?php
declare(strict_types=1);

namespace KurrentDB;

use Generator;
use Spiral\Goridge\Relay;
use Spiral\Goridge\RPC\RPC as GoridgeRPC;
use Spiral\RoadRunner\Environment;

final class Client
{
    private GoridgeRPC $rpc;

    /**
     * @param GoridgeRPC|null $rpc Pass your own, or omit to auto-create from RR ENV.
     */
    public function __construct(?GoridgeRPC $rpc = null)
    {
        /**
         * $address = Environment::fromGlobals()->getRPCAddress();
         * $rpc = new RPC(Relay::create($address));
         */
        $this->rpc = $rpc ?? new GoridgeRPC(
            Relay::create(Environment::fromGlobals()->getRPCAddress())
        );
    }

    /**
     * Read a single chunk (one RPC) — useful for manual paging or small reads.
     *
     * @return array<int, array> Raw event payloads (same shape as Go's ResolvedEvent JSON)
     */
    public function readStream(
        string                  $stream,
        Direction               $direction = Direction::Forwards,
        StreamPosition|int|null $from = StreamPosition::Start,
        ?int                    $maxEvents = null,
        bool                    $resolveLinkTos = false,
        ?int                    $deadline = null,
        bool                    $requiresLeader = false
    ): array
    {
        $payload = [
            'stream' => $stream,
            'direction' => $direction->value,
            'maxEvents' => $maxEvents,
            'from' => $from,
            'resolveLinkTos' => $resolveLinkTos,
            'deadline' => $deadline,
            'requiresLeader' => $requiresLeader,
        ];

        return $this->rpc->call('kurrentdb.ReadStream', $payload);
    }

    /**
     * Read the whole stream lazily as a Generator. Internally pages with readStreamOnce().
     *
     * @return Generator<array> yields each ResolvedEvent (raw array) one-by-one
     */
    public function readStreamPaged(
        string                  $stream,
        Direction               $direction = Direction::Forwards,
        StreamPosition|int|null $from = StreamPosition::Start,
        ?int                    $maxEvents = null,
        bool                    $resolveLinkTos = false,
        ?int                    $deadline = null,
        bool                    $requiresLeader = false,
        int                     $pageSize = 500
    ): Generator
    {
        $cursor = $from;

        while (true) {
            $chunk = $this->readStream(
                stream: $stream,
                direction: $direction,
                from: $cursor,
                maxEvents: $pageSize,
                resolveLinkTos: $resolveLinkTos,
                deadline: $deadline,
                requiresLeader: $requiresLeader,
            );

            if (empty($chunk)) {
                break;
            }

            foreach ($chunk as $resolved) {
                yield $resolved;

                // advance cursor using the last seen EventNumber
                if (isset($resolved['Event']['EventNumber'])) {
                    $cursor = (int)$resolved['Event']['EventNumber'] + 1;
                }
            }

            // If we got fewer than pageSize, assume EOF
            if (count($chunk) < $pageSize) {
                break;
            }
        }
    }

    /**
     * @param string $stream
     * @param EventData[] $events
     * @param StreamState|int $streamState Stream state or revision
     * @param int|null $deadline
     * @param bool $requiresLeader
     * @return WriteResult
     */
    public function appendToStream(
        string          $stream,
        array           $events,
        StreamState|int $streamState = StreamState::Any,
        ?int            $deadline = null,
        bool            $requiresLeader = false
    ): WriteResult
    {
        $payload = [
            'stream' => $stream,
            'events' => $events,
            'deadline' => $deadline,
            'streamState' => $streamState,
            'requiresLeader' => $requiresLeader,
        ];

        $response = $this->rpc->call('kurrentdb.AppendToStream', $payload);
        return WriteResult::fromRpc($response);
    }

}