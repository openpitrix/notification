# Notification
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/KubeSphere/KubeSphere/blob/master/LICENSE)

----

## Introduction
Notification is an enterprise-grade general-purpose high-performance distribute notification system. 

The basic requirements for this system is below:

1.General Purpose

2.Different notification ways

3.Distribute, Asynchronous sending 

4.Notification Address management

It is plugin-driven and designed to support following notification ways:

1.Email

2.Websocket(WIP)

3.Wechat(todo)

4.SMS(todo)

In the future it will provide more functions to support different notification ways.

## Installation:
You can find the details in the [installation documents](doc/installation/allinone.md).
 
 
 
## Architecture Design

![Architecture](doc/images/notification.png)

Notes:

1.Notification provides gRPC and RESTful api for third party call.

2.The Persistence Layer is Mysql.

3.Asynchronous sending Notification, need use MQ to temporarily store notification, using Redis or etcd queue.

 
