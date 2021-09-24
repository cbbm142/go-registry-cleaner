[![codecov](https://codecov.io/gh/cbbm142/go-registry-cleaner/branch/main/graph/badge.svg?token=43UESZJV3S)](https://codecov.io/gh/cbbm142/go-registry-cleaner)

## Go-Registry-Cleaner

This program will scan a supplied docker registry for specific repos and tags.  It will then remove stale images from it to reduce space usage of the registry.

To use supply a config.yml set up for your repoistory and images you want to clean.  Add a .env file with your username and password and then run the program.  

You can also use the included Dockerfile to build and run it as a container, or use the one supplied on Docker Hub.  Here is an example to run it.

`docker run -v $(pwd)/config.yml:/app/config.yml -e username=user -e password=pass -it cbbm142/go-registry-cleaner`

I make no guarantees this works correctly, and you should always use the dryRun option before attempting to delete any images.

### TODO 

- Increase test coverage
- Make use of bearer functions to enable auth with bearer tokens

