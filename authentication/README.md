# Authentication Service
This service has two docker containers:
- prauxy/authentication (./backend): the backend authentication & authorization service.
- prauxy/authentication-web (./frontend): the front end that lets users interact with the backend

## Testing
Since this service relies on two different containers, make sure you build the latest versions first:
**Build backend API**: `docker build -t prauxy/authentication authentication/backend
**Build frontend**: `docker build -t prauxy/authentication-web authentication/frontend

Once you've built the containers, run the following Docker Compose command to test it with Cypress:
```
docker-compose -f authentication/docker-compose.yaml -f compose-helpers/tests-compose.yaml up --exit-code-from cypress
```

## Backend API
This service runs on Go. 

## Frontend UI
The frontend is made in Vue3 and tested with Cypress. 