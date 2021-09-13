
## Go-Registry-Cleaner

This program will scan a supplied docker registry for specific repos and tags.  It will then remove stale images from it to reduce space usage of the registry.

To use supply a config.yml set up for your repoistory and images you want to clean.  Add a .env file with your username and password and then run the program.  

You can also use the included Dockerfile to run it as a container on demand.

I make no guarantees this works correctly, and you should always use the dryRun option before attempting to delete any images.

### T0DO 

- Increase test coverage
- Make use of bearer functions to enable auth with bearer tokens

