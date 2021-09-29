# Code Patterns

An attempt to demonstrate code patterns by solving the same
(synthetic) coding task in different languages.

The basic task is to implement a key-value store.
- The store has finite key capacity. Values are unbounded in size.
- When the store is full and a new key is added, an FIFO policy is used to
  evict existing keys as necessary to keep the size under the limit.
- The store can use different backing implementations (say, memory-backed or
  file-backed).
