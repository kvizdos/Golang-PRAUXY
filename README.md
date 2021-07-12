# Go(lang) PRAUXY
PRAUXY is a tool for developers to manage authentication, create development & production environments, send notifications, and more. Rewritten in Golang using Docker, Go PRAUXY is faster than ever and more customizable. Instead of needing to launch every PRAUXY service, you can now customize your Docker-compose file and only launch what you need. Soon, we'll be launching an interactive Docker-compose creator! 

## Deployment
PRAUXY is now containerized, which makes deployments very simple. If you use Docker Compose, you can just run `docker-compose up` in the root directory to launch all of PRAUXY's services. If you would like to customize the Docker-compose file by removing unneeded services, feel free to do so! 

## Testing
Testing is on a service-level. Please go into each service to see instructions.

## Service
- [x] [Core authentication](/authentication) 
    - [![Authentication Service](https://github.com/kvizdos/Golang-PRAUXY/actions/workflows/authentication.yml/badge.svg)](https://github.com/kvizdos/Golang-PRAUXY/actions/workflows/authentication.yml)
    - Multifactor authentication:
        - [x] TOTP MFA
        - [ ] Hardware keys + backups (FIDO)
    - [ ] OAuth authentication to loop into custom apps
- [ ] Internal PRAUXY service router
    -  This lets PRAUXY be modular and customizable
- [ ] Site deployment
    - [ ] Static sites (w/ automated GitHub updates)
    - [ ] NodeJS deployment
    - [ ] Golang & other language support
    - [ ] Automated Web A/B Testing
- [ ] Reverse Proxy
    - [ ] Native SSL support (LetsEncrypt, custom CA)
- [ ] Certificate Authority (CA) server
    - [ ] Manage certificate access
- [ ] Send & track web notifications
- [ ] PRAUXY One (short URLs w/ advanced tracking)
- [ ] Analytics & User Journey tracking