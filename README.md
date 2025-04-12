# Cache

## SingleCache
* In built Map
* Thread safe
* Uses Read Lock for Read
* Uses RW Lock for writes
* 300k puts in a single core machine in 1 second

## LRUCache
* Thread safe with mutex
* Uses a doubly linked list to track LRU element. Front is most used, and back is least used
