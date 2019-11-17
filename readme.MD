# AWS Hack 'n' Roll - OpenFreePay
  
What is AWS, I make a small joke:  
  
![Joke Title](https://imgur.com/ncXuk16.jpg)
  
  
## Idea
  
* 為了解決更直接跨平台點數兌換的服務需求，因此團隊基於AWS微服務架構來設計此平台應用。
* 本平台提供跨平台點數等值兌換的撮合服務。可供欲兌換點數的使用者族群能有一個既定各家平台的點數幣值標準作為依據，去溝通和交易彼此點數，甚至可彈性調整點數兌換量。
* 使用者能夠上傳目前原有的特定平台點數和欲兌換的平台點數。
  
## Structure
  
![Structure](snapshop.png)  
  
  
## Content
  
#### Enviorment:
* EC2 = AMI: ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20191002 (ami-0cd744adeca97abb1)
* EC2 = AMI: Windows_Server-2019-English-Full-Base-2019.10.09 (ami-0e48df3801c3e668e)
* Docker image: openfreepay/Dockerfile
* Amazon ECS = OpenFreePay-FrontEndPoint cluster
