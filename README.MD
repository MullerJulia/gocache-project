# **GoCache - Thread-Safe In-Memory Cache with LRU Eviction**

## **Overview**

**GoCache** is a thread-safe, in-memory cache system built in Go. It supports basic cache operations (`Set`, `Get`, and `Delete`) and implements a Least Recently Used (LRU) eviction policy to manage cache capacity. This cache is designed to work efficiently in a multi-threaded environment, ensuring data integrity and high performance.

## **Features**

- **Thread-Safe Operations**: Supports concurrent access by multiple goroutines.
- **LRU Eviction**: Automatically evicts the least recently used items when the cache exceeds its capacity.
- **Basic Cache Operations**:
  - `Set(key string, value string)`: Adds or updates a key-value pair in the cache.
  - `Get(key string) (string, bool)`: Retrieves the value associated with a key.
  - `Delete(key string) bool`: Removes a key-value pair from the cache.

## **Getting Started**

### **Prerequisites**

- **Go**: Ensure that Go is installed on your machine. Download it from the [official Go website](https://golang.org/dl/).

### **Installation**

**Clone the Repository:**
```bash
git clone https://github.com/MullerJulia/gocache-project.git
cd gocache-project


## **Running the Example**

Compile and run the example:

```bash
go run main.go

## **Testing**

**GoCache** includes a suite of tests to verify its functionality and thread safety.

### **Running Tests**

Run all tests:

```bash
go test -v ./...

### **Check Test Coverage**

```bash
go test -cover ./...

### **Run Race Tests:**

```bash
go test -race -v ./...
