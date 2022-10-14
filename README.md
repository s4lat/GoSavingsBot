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
  <img width="30%" alt="1" src="https://user-images.githubusercontent.com/24609869/195347856-1c8c2d7a-65f3-4d7a-be14-9eb8d5ca6888.png">
  <img width="30%" alt="2" src="https://user-images.githubusercontent.com/24609869/195347889-a2f1bb1c-2c3e-4c22-9bfc-99e3eb652402.png">
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
The database and logs will be located in `GoSavingsBot/data` folder
# License

GoSavingsBot is distributed under MIT.

