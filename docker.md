### Deployment of Docker image:
---
* Build up images: (at the same level as each of the ```.dockerfile``` in each service.)
```dockerfile
docker build -f [your.dockerfile] -t [dockerusername]/[servicename]:[yourtag]
docker push [dockerusername]/[servicename]:[yourtag]
```