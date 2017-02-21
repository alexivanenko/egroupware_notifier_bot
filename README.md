# EGroupware Telegram Notifier  

## Simple private application for sending notifications to Telegram chat about EGroupware new/modified tasks in GMail inbox    

### Requirements

EGroupware Telegram Notifier version 0.1 requires Go >= 1.7
Installed Telegram Bot: you need bot token and chat ID.   

##### Installation

```sh
$ go get github.com/alexivanenko/egroupware_notifier_bot/...
$ cd egroupware_notifier_bot
$ cp config.sample config.ini
$ mkdir data
$ chmod 0755 data
```

Please update config.ini using appropriate values for Telegram Bot Token and Telegram Chat ID.
Browse to this URL https://developers.google.com/gmail/api/quickstart/go and go through Step1.
During this instruction you will download json file. Move this file to your working directory and rename it ```client_secret.json```.  

```sh
$ make
```

The first time you run the application, it will prompt you to GMail authorize access. Browse to the provided URL in your web browser. Copy the code you're given, paste it into the command-line prompt, and press Enter.
The application checks GMail account every 5 minutes.