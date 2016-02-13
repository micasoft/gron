#gron

Ligthweight job server.

The aim of this project is manager some jobs according priority and time. OK maybe we are inventing the wheel, I know this should be handle by kernel but in some case is very hard to manage this according the complexity and the huge list of task that need to be run at the same time. In my case we had hundred of jobs in cronjob and some of them running at same time. For much I tried to manager this, is difficult to find to the right setup. 

Gron is a small project in golang that you can define a thredshold of maximum of process. At the same time him receive commands and execute them if is available to do this if not we will wait. 

## How use it

### Server
```bash
gron -d -max 20
```

### Client
```bash
gron -c "php worker.php" -p 2
```