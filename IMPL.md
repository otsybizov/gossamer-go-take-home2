# MessageTracker implementation

`MessageTracker` is a mix of queue and random access containers. As such a data structure is required that would allow
to efficiently perform queue and random access operations simultaneously.

A linked list is a natural choice to implement a queue. But considering the requirement to return an array of all
message pointers in the order in which they were received on demand, a linked list would badly impact the performance
in this case because O(n) copy operations, where n is the number of messages, would be required to build the array
to return. So I considered storing message pointers in an array, that would be ready to be returned any moment.

A hash table is a good choice to implement random access operations. The problem with it is that the order of elements
is not preserved.

Considering the above I decided to sacrifice space in order to use advantages of both array and hash table and duplicate
message pointers in both container types. This looks to me a reasonable tradeoff because pointers are stored, not
instances of message structure.

## Add

Add operation appends a message pointer to the end of the array and inserts it into the map. When the maximum number of
messages is reached, Add operation in addition removes the first element from the array and the element with the
same ID from the map.

## Delete

Delete operation looks up and deletes the message with the specified ID from the array and from the map. This is not
a very efficient way because of the decision to keep the array up to date to return all messages immediately. Perhaps
this can be changed based on an assumption how frequently `Messages` would be called, e.g. store null pointers in
the slots where the messages were deleted and remove them if and when the slots are required for new elements or when
all messages are requested.

## Message (get one)

Message operation looks up the element with the specified ID in the map. If the element found then it is returned,
otherwise `message not found` error returned.

## Messages (get all)

Returns the copy of the array.




