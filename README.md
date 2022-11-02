# Ankorstore CLI plugin for Aiven Services

This work-in-progress repo contains a Go implementation of a future plugin for the Ankorstore CLI utility, letting users manage Aiven services.

## Prerequisites
In order for this to work you need an account on an Aiven subscription. 
Create an API token for your account (plugin doesn't currently support OAuth2 authentication) and store it to a local env-var `AIVEN_API_TOKEN`

## Usage
The plugin supports the following commands

### Listing projects
```ankor-aiven-plugin projects```
This will list all projects your account has access to, alongside with the role your account holds.

### Listing services
```ankor-aiven-plugin services [projectName]```
This will list all the services associated with the project that you provide as an argument. For each service you will also get the current state of it.

### Stopping a service
```ankor-aiven-plugin stop [projectName] [serviceName]```
This will stop the service given as an argument, which is part of the project provided as first argument.

### Starting a service
```ankor-aiven-plugin start [projectName] [serviceName]```
This will startt the service given as an argument, which is part of the project provided as first argument.
