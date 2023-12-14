# Notes

## The RCON protocol

The RCON protocol is documented by Valve here: <https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Requests_and_Responses>.

Sadly, this documentation is not exactly a specification.
Due to this, games that use the RCON protocol may implement it slightly differently.

### Squad RCON

#### Multi-packet responses

If the response is too large, it is split over multiple packets.
Sadly, packets do not contain an indicator on whether any further packets will be sent.

The documentation recommends that after sending a command, to immediately
[send an empty SERVERDATA_RESPONSE_VALUE](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Multiple-packet_Responses).

> One common workaround is for the client to send an empty `SERVERDATA_RESPONSE_VALUE` packet after
> every `SERVERDATA_EXECCOMMAND` request.
> Rather than throwing out the erroneous request, SRCDS mirrors it back to the client,
> followed by another RESPONSE_VALUE packet containing 0x0000 0001 0000 0000 in the packet body
> field.
> Because SRCDS always responds to requests in the order it receives them,
> receiving a response packet with an empty packet body guarantees that all of the meaningful
> response packets have already been received.
> Then, the response bodies can simply be concatenated to build the full response.

However, while Squad reflects the `SERVERDATA_RESPONSE_VALUE`, it does not continue following the documented behavior.
Rather, it does **not** send the `RESPONSE_VALUE` packet with body `0x0000 0001 0000 000` but
it sends the confirmation packet twice followed by `00 01 00 00` (hex).

Instead, what we do is; after a command (`C`), we send another command where we do not care for
the answer.
We give this _confirmation command_ the ID of `C.Id` incremented by one.
By ensuring the ID of `C` is always even, we know when receiving a packet with an odd ID that
the original command is complete.
