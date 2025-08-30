<?php

use KurrentDB\Client;
use KurrentDB\EventData;
use Nyholm\Psr7\Factory\Psr17Factory;
use Nyholm\Psr7\Response;
use Ramsey\Uuid\Uuid;
use Spiral\RoadRunner\Http\PSR7Worker;
use Spiral\RoadRunner\Worker;


// Create new RoadRunner worker from global environment
$worker = Worker::create();

// Create common PSR-17 HTTP factory
$factory = new Psr17Factory();

$psr7 = new PSR7Worker($worker, $factory, $factory, $factory);

while (true) {
    try {
        $request = $psr7->waitRequest();
        if ($request === null) {
            break;
        }
    } catch (Exception $e) {
        $psr7->respond(new Response(400));
        continue;
    }

    try {
        /*
        $address = Environment::fromGlobals()->getRPCAddress();
        $rpc = new RPC(Relay::create($address));
        $request = new KurrentDB\ReadStreamRequest('MyDocument');
        $res = $rpc->call('kurrentdb.ReadStream', $request);
*/
        $client = new Client();

        $writeResult = $client->appendToStream("MyDocument",
            [
                new EventData(
                    Uuid::uuid7()->toString(),
                    "EditorChanged",
                    ["editor" => "Edvin Syse"]
                ),
            ],
        );

        $psr7->respond(new Response(200, ['Content-Type' => 'application/json'],
            json_encode($writeResult)
        ));

        $res = $client->readStream("MyDocument");
        $psr7->respond(new Response(200, ['Content-Type' => 'application/json'],
            json_encode($res)
        ));

    } catch (Exception $e) {
        // In case of any exceptions in the application code, you should handle
        // them and inform the client about the presence of a server error.
        //
        // Reply by the 500 Internal Server Error response
        $psr7->respond(new Response(500, [], 'Something Went Wrong: ' . $e->getMessage()));

        // Additionally, we can inform the RoadRunner that the processing
        // of the request failed. Use error instead of response to indicate
        // worker error, do not use both.
        // $psr7->getWorker()->error((string)$e);
    }
}