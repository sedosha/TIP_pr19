sequenceDiagram
    participant Client as Клиент
    participant Tasks as Tasks Service<br/>(порт 8090)
    participant Auth as Auth Service<br/>(порт 8089)
    participant Storage as In-memory Storage

    Note over Client,Tasks: Сценарий 1: Получение токена
    
    Client->>Auth: POST /v1/auth/login
    Auth-->>Client: 200 OK + access_token
    
    Note over Client,Tasks: Сценарий 2: Создание задачи
    
    Client->>Tasks: POST /v1/tasks + Authorization: Bearer token
    Tasks->>Auth: GET /v1/auth/verify + Authorization
    Auth-->>Tasks: 200 OK (valid=true)
    Tasks->>Storage: Сохранить задачу
    Storage-->>Tasks: task_id
    Tasks-->>Client: 201 Created + task data
    
    Note over Client,Tasks: Сценарий 3: Неверный токен
    
    Client->>Tasks: GET /v1/tasks + Authorization: Bearer wrong-token
    Tasks->>Auth: GET /v1/auth/verify + Authorization
    Auth-->>Tasks: 401 Unauthorized
    Tasks-->>Client: 401 Unauthorized
    
    Note over Client,Tasks: Сценарий 4: Auth недоступен
    
    Client->>Tasks: GET /v1/tasks + Authorization: Bearer token
    Tasks->>Auth: GET /v1/auth/verify (timeout 3s)
    Auth--xTasks: Service Unavailable
    Tasks-->>Client: 503 Service Unavailable
