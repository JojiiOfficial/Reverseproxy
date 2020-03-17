# Reverseproxy
A high performance reverse proxy written in go. It can replace reverseproxies like apache or nginx if you don't need special features. It is pretty simple to setup, lightweight and very fast. In some cases it is faster than apache/nginx!

# Installation
You can compile the binary by using `go build -o main`. This will create a file called `main`.<br>
You can set `-debug` and `-config` if you want to view more output or specify a custom config file. By default the config file is stored at /etc/reverseproxy/config.toml

### Docker
The images are hosted at [Dockerhub](https://hub.docker.com/r/jojii/reverseproxy)
To Install run following command (you can/should adjust the volume path)
```bash
docker run --name revproxy --rm -v `pwd`/config:/app/config jojii/reverseproxy
```

## Concept/Idea
- You have one configfile in which you have to define your routefiles and interfaces
- You can have nroutes stored in separate files in the `./config/routes` directory
- One route represents one (sub)domain/host
- One route can listen on n ports/interfaces which you have to specify in the config first
- Http redirect interfaces can't be used as interface for locations

# Configuration
### Example
Config.toml:
```toml
# Specify your routes
RouteFiles = ["./config/routes/route1.toml"]
[Server]
  MaxHeaderSize = "16KB"
  ReadTimeout = "10s"
  WriteTimeout = "10s"
  
# Setup port 80 as auto http redirect (to https)
[[ListenAddresses]]
  Address = ":80"
  SSL = false
  Task = "httpredirect"
  [ListenAddresses.TaskData]
    [ListenAddresses.TaskData.Redirect]
      HTTPCode = 301

# Use 443 using SSL 
[[ListenAddresses]]
  Address = ":443"
  SSL = true

```
route1.toml:
```toml
ServerNames = ["yourDomain.xyz"]
Interfaces = [":80", ":443"]

# Your ssl stuff
[SSL]
  Key = "./certs/key.pem"
  Cert = "./certs/cert.pem"

[[Location]]
  # Location to match for this route
  Location = "/"
  # Destination (must be a URL to somewhere)
  Destination = "http://127.0.0.1:81/"
  # Is regex in Location
  Regex = false  
  
[[Location]]
  Location = "/hidden/secret/stuff"
  Destination = "http://127.0.0.1:81/admin/"
  # Only allow localhost and 192.168.1.1/24 to access this location
  Deny = "all"
  Allow = ["127.0.0.1", "192.168.1.1/24"]
```

## Important
- You <b>must</b> specify every interface you use in routes in the config exact the same way!
- You should put the root location (/) at the end of your locations. The priority is from top to bottom
