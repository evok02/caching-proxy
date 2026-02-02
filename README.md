# Cacher: proxy server 

*Hand-made* version of caching proxy server written in go and containerized using Docker.

## Motivaion
The goals of this project was to experiment with multiple instances of the app using containerization. Understanding the concepts behind the images and containers, as well as basic understanding of how they work behind the scene.
Among interesting topics for me was namespaces, chroot and mounts as a primary source of "isolated" environment. Also i gained first experience configuring servers, in this project i used Caddy as a load-balancer.

## How does it work?
The basic scheme of how the setup works behind the scene you can check below.

<img width="975" height="777" alt="image" src="https://github.com/user-attachments/assets/bb08ee38-7e11-4ca1-a852-e887ef8befaf" />

## Quick start?

- Obviously you need [docker engine](https://www.docker.com/get-started/)

- Clone the repo on your machine:

  ```
  git clone https://github.com/evok02/cacher
  ```
- Creatig bridge:

  ```
  make create-bridge
  ```
- Run builder command:

  ```
  make run-proxy-container
  ```

## Usage

You can configure instances of the containers using **./configs/main.yml**. Below you can see configuration by default:
```
proxy:
  port: "80"
  host: "0.0.0.0"
  # url prefix for proxy
  target: "https://github.com/"  
  # log files path configuration 
  log: 
    err_log: "./log/error.log" 
    info_log: "./log/info.log"
  db: 
    port: "6379"
    # name used to resolve address of the instance within the network
    name: "cache" 

reverse_proxy:
  expose: "8080"
```
## Contributing

I would love your feedback. Contribute by opening pull requests and issues.
