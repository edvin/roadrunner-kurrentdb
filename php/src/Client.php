<?php
declare(strict_types=1);

namespace KurrentDB;

use Generator;
use Spiral\Goridge\Relay;
use Spiral\Goridge\RPC\Codec\MsgpackCodec;
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
        $this->rpc = $rpc ?? new GoridgeRPC(
            Relay::create(Environment::fromGlobals()->getRPCAddress()),
            new MsgpackCodec()
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
            'Stream' => $stream,
            'Direction' => $direction->value,
            'MaxEvents' => $maxEvents,
            'From' => $from instanceof StreamPosition
                ? ['Kind' => $from->value]
                : ['Kind' => 'Index', 'Index' => $from],
            'ResolveLinkTos' => $resolveLinkTos,
            'Deadline' => $deadline,
            'RequiresLeader' => $requiresLeader,
        ];

        return $this->rpc->call('kurrentdb.ReadStream', $payload);
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
            'Stream' => $stream,
            'Events' => array_map(fn(EventData $e) => $e->toArray(), $events),
            'Deadline' => $deadline,
            'StreamState' => $streamState instanceof StreamState
                ? ['Kind' => $streamState->value]
                : ['Kind' => 'Revision', 'Revision' => $streamState],
            'RequiresLeader' => $requiresLeader,
        ];

        $response = $this->rpc->call('kurrentdb.AppendToStream', $payload);
        return WriteResult::fromRpc($response);
    }

    /**
     * @param string $stream
     * @param StreamState|int $streamState Stream state or revision
     * @param int|null $deadline
     * @param bool $requiresLeader
     * @return DeleteResult
     */
    public function deleteStream(
        string          $stream,
        StreamState|int $streamState = StreamState::Any,
        ?int            $deadline = null,
        bool            $requiresLeader = false
    ): DeleteResult
    {
        $payload = [
            'Stream' => $stream,
            'Deadline' => $deadline,
            'StreamState' => $streamState instanceof StreamState
                ? ['Kind' => $streamState->value]
                : ['Kind' => 'Revision', 'Revision' => $streamState],
            'RequiresLeader' => $requiresLeader,
        ];

        $response = $this->rpc->call('kurrentdb.DeleteStream', $payload);
        return DeleteResult::fromRpc($response);
    }

}