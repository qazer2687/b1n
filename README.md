## b1n
a minimal file drop service

---

### about
The project is written in go with zero third-party dependencies. Uploads are streamed to disk, which keeps memory low and bounded, and hashed with SHA-256. The hex digest of the hash becomes the file identifier, which is computationally infeasible to enumerate. The service is configured through environment variables.

### endpoints

`POST /upload` - upload a file

`GET /{id}` - download a file

`GET /` - web interface