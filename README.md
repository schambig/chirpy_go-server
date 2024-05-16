Building a Go Web Server from scratch.

A web server is just a computer that serves data over a network, typically the Internet.

Servers run software that listens for incoming requests from clients. When a request is received, the server responds with the requested data.

```mermaid
graph TD
    A[Client] -->|HTTP Request| B[Server]
    B -->|Route| C[Router]
    C -->|Controller| D[Controller]
    
    D -->|Process JSON| E[Service]
    E -->|Data| F[Storage]
    E -->|Auth Request| G[Auth Service]
    
    subgraph Storage
        F1[Database]
        F2[Cache]
    end
    F1 --> F
    F2 --> F
    
    subgraph Auth
        G1[Login]
        G2[Token]
    end
    G1 --> G
    G2 --> G
    
    G -->|Verify| H[Authentication]
    G -->|Check| I[Authorization]
    H -->|Auth Result| E
    I -->|Auth Result| E
    
    E -->|Trigger| J[Webhooks]
    J -->|Send Data| K[External Service]
```
Brief explanation of the key points from the graph above:
- **Client**: Initiates the HTTP request.
- **Server**: Handles the incoming HTTP request.
- **Router**: Directs the request to the appropriate controller.
- **Controller**: Processes the request and interacts with the service layer.
- **Service**: Manages business logic, interacts with storage, and handles authentication/authorization.
- **Storage**: Consists of a database and cache for storing data.
- **Auth Service**:
  - Login: Manages user login.
  - Token: Handles token creation and verification.
- **Authentication**: Verifies user credentials.
- **Authorization**: Checks user permissions.
- **Webhooks**: Manages webhooks and sends data to external services.
- **External Service**: Receives data from webhooks.
