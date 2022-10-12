# GoSavingsBot

```bash
git clone https://github.com/s4lat/GoSavingsBot
```

* [Overview](#overview)
* [Deploy](#deploy)
* [License](#license)

# Overview
GoSavingsBots is a telegram bot that helps you keep track of your spending.  
You can deploy it for yourself or use an already deployed one - [@GoSavingsBot](https://t.me/GoSavingsBot)
<p float="left">
  <img width="200" alt="1" src="https://user-images.githubusercontent.com/24609869/195345424-13e1348c-3087-4b2d-9975-fd9ffd265132.png">
  <img width="200" alt="2" src="https://user-images.githubusercontent.com/24609869/195344595-d3fd7bab-b5f3-4164-8399-bb493f5fc7bb.png">
  <img width="200" alt="3" src="https://user-images.githubusercontent.com/24609869/195344608-7686df09-c9bc-4e77-bf3e-c0236ef75db5.png">
  <img width="200" alt="4" src="https://user-images.githubusercontent.com/24609869/195344634-4b1b52d1-5baf-4b89-a5b6-a5052636cd0d.png">
 </p>

**GoSavingsBot** features:
* Adding expenses with the following message format: `<cost> <name>`
* Using **time zones** to avoid time inaccuracies
* Display statistics for the **day**, **month** and **year**
* 🇬🇧 **English** and 🇷🇺 **Russian** language support
* Possibility to export the list of spends to **csv/excel** file
* Ability for the user to **delete all his data** from the database

# Deploy
Put bot token in the `BOT_TOKEN` field in [docker-compose.yml](docker-compose.yml) and run:
```
docker-compose up -d
```
The database with user data will be located in `GoSavingsBot/data` folder
# License

GoSavingsBot is distributed under MIT.

