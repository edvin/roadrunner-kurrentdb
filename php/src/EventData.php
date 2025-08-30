<?php

namespace KurrentDB;

use JsonSerializable;

final class EventData implements JsonSerializable
{
    public function __construct(
        public string      $eventId,
        public string      $eventType,
        public             $data,
        public             $metadata = null,
        public ContentType $contentType = ContentType::Json,
    )
    {
    }

    public function jsonSerialize(): array
    {
        return [
            'EventID' => $this->eventId,
            'EventType' => $this->eventType,
            'ContentType' => $this->contentType->value,
            'Data' => base64_encode(json_encode($this->data)),
            'Metadata' => $this->metadata !== null ? base64_encode(json_encode($this->metadata)) : '',
        ];
    }
}