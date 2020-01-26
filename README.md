
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


### Starting Data
4