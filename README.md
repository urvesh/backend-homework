
## Local Setup for Dev

Requires Homebrew, Go and MongoDB

Brew changed where mongo is located, so now we have to get it 
by doing the following:
Also [Link](https://github.com/mongodb/homebrew-brew)
```bash
brew tap mongodb/brew
brew install mongodb-community
sudo brew services start mongodb-community
```

Mongo conf location:
`/usr/local/etc/mongod.conf`



## EC2 Setup Notes

https://docs.mongodb.com/manual/tutorial/install-mongodb-on-amazon/
https://hackernoon.com/deploying-a-go-application-on-aws-ec2-76390c09c2c5


```bash
sudo yum install -y git
sudo yum install -y golang
```

To run process in background:
https://linuxhandbook.com/run-process-background/
```bash
make run &
Ctrl + C

nohup make run &
```

To kill background process:
```bash
jobs
fg
Ctrl + C

ps -aux | grep backend
kill -9 pid
```

## Starting Data
#### Users

| _id                      | name     |
| ------------------------ | -------- |
| 5e2e39ee290f5a56ffda9ed5 | Jennifer |
| 5e2e39ee290f5a56ffda9ed6 | Bob      |
| 5e2e39ee290f5a56ffda9ed7 | Susan    |
| 5e2e39ee290f5a56ffda9ed8 | Michael  |
| 5e2e39ee290f5a56ffda9ed9 | Alexis   |
| 5e2e39ee290f5a56ffda9eda | Andrew   |

#### Ratings
| fromUserId               | toUserId                 | type |
| ------------------------ | ------------------------ | ---- |
| 5e2e39ee290f5a56ffda9ed5 | 5e2e39ee290f5a56ffda9ed8 | LIKE |
| 5e2e39ee290f5a56ffda9ed6 | 5e2e39ee290f5a56ffda9ed5 | LIKE |
| 5e2e39ee290f5a56ffda9ed6 | 5e2e39ee290f5a56ffda9ed7 | LIKE |
| 5e2e39ee290f5a56ffda9ed6 | 5e2e39ee290f5a56ffda9ed8 | LIKE |
| 5e2e39ee290f5a56ffda9ed7 | 5e2e39ee290f5a56ffda9ed5 | LIKE |
| 5e2e39ee290f5a56ffda9ed8 | 5e2e39ee290f5a56ffda9ed5 | LIKE |
| 5e2e39ee290f5a56ffda9ed9 | 5e2e39ee290f5a56ffda9ed5 | LIKE |
| 5e2e39ee290f5a56ffda9ed9 | 5e2e39ee290f5a56ffda9ed8 | LIKE |
| 5e2e39ee290f5a56ffda9ed9 | 5e2e39ee290f5a56ffda9eda | LIKE |
| 5e2e39ee290f5a56ffda9eda | 5e2e39ee290f5a56ffda9ed5 | LIKE |

Notes:
* Jennifer likes Michael (match)
* Bob likes Jennifer, Susan, and Michael
* Susan likes Jennifer
* Michael likes Jennifer
* Alexis likes Jennifer, Michael and Andrew
* Andrew likes Jennifer

