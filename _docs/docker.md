
# Build & Deploy RestAPI Image

```bash
docker build -t api:multistage -f Dockerfile.multistage.api .
```

# Tag Image

```bash
docker image tag api:multistage api:v1.0
```

# Run Image as Container

```bash
docker run -p 8080:8080 api:v1.0 --name api
```

