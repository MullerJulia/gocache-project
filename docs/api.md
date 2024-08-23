# GoCache Documentation

## Overview

This document provides an overview of the GoCache implementation, including concurrency mechanisms, the LRU eviction strategy, code documentation, and usage examples.

## Concurrency Mechanisms

### RWMutex

- **Description**: The cache uses `sync.RWMutex` to manage concurrent access to the cache.
- **Purpose**: Allows multiple goroutines to read the cache concurrently while ensuring exclusive access for write operations and eviction processes.

### Set Method

- **Lock Type**: Write lock
- **Purpose**: Acquires a write lock to safely update or insert new key-value pairs into the cache. This prevents data races and ensures that no other goroutine can read or write to the cache while the operation is in progress.
- **Usage**:
    ```go
    cache := gocache.NewGoCache(3)
    cache.Set("key1", "value1")
    cache.Set("key2", "value2")
    ```

### Get Method

- **Lock Type**: Read lock initially, upgraded to write lock if necessary
- **Purpose**: Starts with a read lock to allow concurrent reads of the cache. If the key is found and needs to be moved to the front of the LRU list, the lock is upgraded to a write lock to safely perform this operation.
- **Usage**:
    ```go
    cache := gocache.NewGoCache(3)
    cache.Set("key1", "value1")
    
    value, found := cache.Get("key1")
    if found {
        fmt.Println("Value:", value)
    } else {
        fmt.Println("Key not found")
    }
    ```

### Delete Method

- **Lock Type**: Write lock
- **Purpose**: Acquires a write lock to safely remove key-value pairs from the cache. This ensures that no other goroutine can access or modify the cache while the deletion is in progress.
- **Usage**:
    ```go
    cache := gocache.NewGoCache(3)
    cache.Set("key1", "value1")
    
    deleted := cache.Delete("key1")
    if deleted {
        fmt.Println("Key deleted")
    } else {
        fmt.Println("Key not found")
    }
    ```

### Evict Method

- **Purpose**: Internal method that removes the least recently used item from the cache when the capacity is exceeded. This ensures that the cache maintains its size within the specified limit.
- **Usage**: The `evict` method is automatically called by the `Set` method when the cache reaches its capacity. It is not meant to be called directly by users of the cache.

## LRU Eviction Strategy

### List and Map

- **Description**: The cache uses a combination of a doubly-linked list (`container/list`) and a map to implement the Least Recently Used (LRU) eviction policy.
- **Purpose**:
  - **Map**: Provides constant-time lookups for cache entries.
  - **List**: Maintains the order of access, allowing efficient eviction of the least recently used item when the cache exceeds its capacity.

### **Build and Run the Example**

Open a terminal and navigate to the directory containing `main.go` file. Run the following command to build and execute the example:

```bash
go run main.go

### **Expected Output**

When you run the example, you should see output indicating the actions performed on the cache, such as:

```plaintext
Key a: apple
Key b was evicted
Key c was deleted