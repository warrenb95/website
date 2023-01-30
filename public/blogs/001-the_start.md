# The start

## Simple bud server

I stumbled across this [go](https://go.dev/) project called [bud](https://github.com/livebud/bud) that looked very interesting having worked with other languages before that had their own full stack solutions like [django](https://www.djangoproject.com/).
Bud looked like a cool project and I'd been planning on building a fullstack website to host blogs and logs of things that I'm learning about or solutions to software engineering related issues; essential a brain dump of everything relating to programming, architecture, ci/cd and cloud.

I won't go into the details of how to implement a bud server in this blog, but should release something in the next few months explaining what I've done 😄.
I'd recommend you check out the [github repo](https://github.com/livebud/bud) for the video tutorial and links to the documentation. 

My bud application at the start is very simple and consists of a single page with a nav bar, title and small paragraph explaining that this is my website and I'm planning to post/update it month.

I also used [bootstrap](https://getbootstrap.com/) to style the website because I'm very lazy and am not a front end developer 😄

## Docker file

To be able to run the app anywhere in the simplest way possible I decided to containerize it and run it on [docker](https://www.docker.com/).
This would allow me to run the app anywhere docker is installed, I won't have to worry about installing programming languages, dependencies or anything else; just install docker and run the container.
I was lucky enought to find a [docker file](https://github.com/livebud/bud/blob/main/contributing/Dockerfile) on the bud repo and I used this as my base, I then added a couple extra commands to expose port `:3000` and run the application.


I update the image when I saw that the image size was `1.82GB`.

![previous docker image size of 1.82GB](/images/docker_size_before.PNG)

I created a multi build image, so I'd have the first part of the dockerfile pull in all the dependencies and create the executable to run the website.
I then copied this executable onto a different image that didn't have the bloat of the dependencies; the executable packages everything it needs to run the website so we don't need anything else on the container.
Using this multi build method I was able to get the image size down to `173MB`, that's 10 times smaller, not a bad reduction in size.

![newer docker image size of 177MB](/images/docker_size_after.PNG)

Below you can see the dockerfile and how I've used a `build` image to build the project and then I copied the executable into a smaller base image `debian:latest`.
I could have used a smaller image to run the project but I had to use debian because this the OS that I built the executable on.

```dockerfile
ARG NODE_VERSION=19.0.0
FROM node:${NODE_VERSION}-slim as build

ARG GO_VERSION=1.18.3
ARG BUD_VERSION=main

RUN node -v

# Install basic dependencies
RUN apt-get -qq update \
  && apt-get -qq -y install curl git make gcc g++ \
  && rm -rf /var/lib/apt/lists/*

# Install Go
RUN curl -L --output - https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -xz -C /usr/local --strip-components 1
RUN go version
ENV PATH "/root/go/bin:${PATH}"

# Install Bud
RUN git clone https://github.com/livebud/bud /bud
WORKDIR /bud
RUN git checkout $BUD_VERSION
RUN make install
RUN go install .
RUN bud version

# Build your project for production
WORKDIR /builder
COPY . . 
RUN go mod download
RUN npm install
RUN bud build
RUN ls -l
RUN ls -l bud/

FROM debian:latest
WORKDIR /app
COPY --from=build /builder/bud/app .
RUN ls -l
EXPOSE 3000

# Run the app
ENTRYPOINT [ "./app" , "--log", "debug", "--listen", "0.0.0.0:3000" ]
```

## AWS setup

### EC2

I went for the quickest/simplest approach and that was to just create an ec2 instance on AWS and manually set up the server. If you don't know what an [ec2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EC2_GetStarted.html?icmpid=docs_ec2_console) instance is it's a virtual machine running on AWS infrastructure.
This allows me to create a virtual machine running linux and then install all the software I need onto it and serve my website.
The first thing I installed was [docker](https://www.docker.com/) which would allow me to run the container image that we previously created, and by using docker to host the app on the ec2 instance I don't have to install any other software, I just need to run a command like this `docker run -it -d --name website {my_image}` and docker will handle the rest.

Now that the website is running on a machine I need someway to access this 🤔. How?
To access the ec2 instance I could create a simple ec2 [security group](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/security-group-rules.html) to allow all HTTP traffic on any port to connect to the ec2 instance and therefore connect to the website. This would work but then there comes the issue of scaling the website.

### Load Balancing

So what I should have done next is create an auto scaling group, I could then have built a template by snapshoting the previously created ec2. The auto scaling group could then spin up new instances if say the average CPU load was >80%. I didn't do this though, but I'll defineteley add this in at some point soonish.
I set up an application load balancer ALB, this would then allow me to configure how the traffic was routed depending on the HTTP port for example. I created a listener on the ALB that would listen for HTTPS traffic on port 443, I also had to set up TLS (transport layer security) certificates using [AWS Certificate Manager](https://docs.aws.amazon.com/acm/index.html). This listener would then forward all traffic to a target group that will route incoming traffic to port `:3000` on my ec2 instance (this will later be my auto scaling group), the port `:3000` is just the port that `BUD` defaults to so I just kept it for now.

I could now access my website using this DNS name `BlogALB-827866070.eu-west-2.elb.amazonaws.com` (this worked as of 12 Jan 2023), this URL is crap and no one (including me) will never remember this. So how can we update this to have a 'normal' looking URL?

## Route 53

[Route 53](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/Welcome.html) is an AWS service that allows you to register domain names, create a hosted zone and then route traffic on your domain to other AWS services of any IP address of your choice. So I registered the domain name `warrenb95.com` and boom 💥 I know have my own domain name and with some Route 53 magic I was able to route traffic from the domain DNS to the ALB.

This is pretty much the current state of the website. If you wanna learn more about how to do this then I'd highly recommend this [udemy course](https://www.udemy.com/course/aws-certified-developer-associate-dva-c01/), this is what I'm currently going through and learning lots about AWS and how to use the services to host my own website.


## Conclusion

To wrap this up what I've done so far is:
1. Create a simple webserver with golang and bud
2. Created a dockerfile to containerize the app so that I can run it basically anywhere
3. Set up an ec2 instance and pulled on the docker container and started the server
4. Create an ALB and routed traffic to my ec2 instance
5. Registered my domain `warrenb95.com` and then created DNS records to route the traffic to my ALB

## What's next?

I'm planning on writing blogs on each of these topics to explain then in more detail and give more examples of how I've used the different technologies.

Thank you so much for checking out the blog post and the website! 😄 

