# PRAUXY Authentication
PRAUXY Authentication is primarily used by core services, however OAuth support will let developers loop simple, yet secure, authorizations into their own apps. 

## Deployment
1. Build Docker container: `docker build -t prauxy/authentication`
2. Use the docker_compose file from [../docker-compose.yaml](../docker-compose.yaml), or create your own.

## Testing
It is recommended to use a testing Docker environment to do e2e tests. You can use the `./test.sh` file in [../authentication-tests](../authentication-tests) folder to easily tear that up and down. Make sure you build this again whenever you update before rerunning tests.