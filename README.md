## b1n
a minimal file drop service

---

**EXPERIMENTAL** - This project is provided as-is without any warranties, express or implied. Use at your own risk.

### about
The project is written in go with zero third-party dependencies. Uploads are streamed to disk, which keeps memory low and bounded, and hashed with SHA-256. The hex digest of the hash becomes the file identifier, which is computationally infeasible to enumerate. The service is configured through environment variables.

### installation

b1n supports docker, there is a `compose.yaml` included in the repository and packages are automatically generated from the main branch and hosted on GHCR at `ghcr.io/qazer2687/b1n:latest`.

### endpoints

`POST /upload` - upload a file

`GET /{id}` - download a file

`GET /` - web interface