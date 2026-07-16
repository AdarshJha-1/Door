# Proxy Server v1 Notes

### Things I learned

* A reverse proxy does **not** simply pass the client's request to the backend. It receives the request as an HTTP server, creates a new request as an HTTP client, sends it to the backend, gets the response, and streams that response back to the client.

* There are actually **two HTTP connections**:

  * Client ↔ Proxy
  * Proxy ↔ Backend

* The proxy should stay as transparent as possible. If the backend returns `404`, `500`, or any other status code, the proxy should usually return the same status code instead of replacing it.

* If the proxy cannot reach the backend, the correct response is generally **502 Bad Gateway**.

* Streaming the response (`io.Copy`) is better than reading the entire body into memory first. It uses less memory and works better for large responses.

* The `Host` header needs special handling. Sometimes the proxy preserves the original host, and sometimes it rewrites it to the backend's host. This becomes important when working with virtual hosts and multiple domains.

---

### Things to remember while creating the backend request

* HTTP method (`GET`, `POST`, etc.)
* URL path
* Query parameters
* Request body
* Request headers
* `Host` header (may need special handling)


---

### Request Flow

```text
Client
   |
   | HTTP Request
   v
Proxy Server
   |
   | Create a new request
   v
Backend Server
   |
   | HTTP Response
   v
Proxy Server
   |
   | Copy status
   | Copy headers
   | Stream body
   v
Client
```
