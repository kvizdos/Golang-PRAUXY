# Compose Helpers
Some Docker Compose things are reused a lot throughout the project, so this directory stores a shared store of them. 

## Example Usage
Below is a docker-compose command that'll spin up the Authentication service & test it. 
```
docker-compose -f authentication/docker-compose.yaml -f compose-helpers/tests-compose.yaml up --exit-code-from cypress
```


## List of helpers:
- cy-display.yaml: if you run a frontend test suite in Docker and would like to get Crypress's nice app to run, use this file.
    - set-x11-host.sh: if you use cy-display.yaml, you need to set your x11-host.
- tests-compose.yaml: runs Cypress