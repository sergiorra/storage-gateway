<div align="center">
  <h1>Storage Gateway</h1>
</div>

## ⚙️ Usage and examples

### Run project

```
make docker-up  //Cleans up containers and then starts up storage Gateway and its dependecies
```

### Request examples

```
PUT localhost:3000/object/weg231 (insert file in the request body)
GET localhost:3000/object/weg231
```


## 📜 Information

### Project Architecture
A quick Medium reading about Hexagonal Architecture in Go: [Hexagonal Architecture](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3)

### TODO
- Add testing
- Expose metrics and run Grafana on Docker container
- Etc...
