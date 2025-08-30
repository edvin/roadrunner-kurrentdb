<?php

namespace KurrentDB;

use Ramsey\Uuid\Uuid;

final class EventData
{
    public function __construct(
        public string            $eventId,          // canonical UUID string
        public string            $eventType,
        public array|string      $data,
        public array|string|null $metadata = null,
        public ContentType       $contentType = ContentType::Json,
    )
    {
    }

    public function toArray(): array
    {
        $eventIdBin = Uuid::fromString($this->eventId)->getBytes(); // 16 bytes

        $data = $this->data;
        if ($this->contentType === ContentType::Json && is_array($data)) {
            $data = json_encode($data, JSON_UNESCAPED_SLASHES);
        } else {
            $data = (string)$data; // raw/binary ok
        }

        $meta = $this->metadata;
        if ($meta !== null) {
            $meta = ($this->contentType === ContentType::Json && is_array($meta))
                ? json_encode($meta, JSON_UNESCAPED_SLASHES)
                : (string)$meta;
        }

        return [
            'EventID' => $eventIdBin,                 // 16-byte bin → UUID ok
            'EventType' => $this->eventType,
            'ContentType' => $this->contentType->value,   // string alias in Go
            'Data' => $data,                       // []byte in Go
            'Metadata' => $meta,                       // []byte or nil
        ];
    }
}