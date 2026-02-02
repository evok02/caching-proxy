# Cacher: proxy server 

*Hand-made* version of caching proxy server written in go and containerized using Docker.

## Motivaion
The goals of this project was to experiment with multiple instances of the app using containerization. Understanding the concepts behind the images and containers, as well as basic understanding of how they work behind the scene.
Among interesting topics for me was namespaces, chroot and mounts as a primary source of "isolated" environment. Also i gained first experience configuring servers, in this project i used Caddy as a load-balancer.

## How does it work?
The basic scheme of how the setup works behind the scene you can check below.

<img width="975" height="777" alt="image" src="https://github.com/user-attachments/assets/bb08ee38-7e11-4ca1-a852-e887ef8befaf" />

## How to run it?

- Obviously you need [docker engine](https://www.docker.com/get-started/)
- Clone the repo on your machine:
  ```
  git clone https://github.com/evok02/caching-proxy
  ```
- Creatig bridge:
  ```
  make create-bridge
  ```
- Run builder/run  command:
  ```
  make run-proxy-container
  ```
